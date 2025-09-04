import { APIGatewayProxyEvent } from 'aws-lambda';
import { Octokit } from '@octokit/rest';
import { z } from 'zod';

// Schemas
import { UserSchema } from '../../models/user.js';
import { SolutionSchema } from '../../models/solution.js';
import { mapProblemLanguageToExtension } from '../../utils/mapProblemLanguage.js';

const commitFormSchema = z.object({
    user: UserSchema,
    solution: SolutionSchema,
    user_access_token: z.string()
});

export const handler = async function (event: APIGatewayProxyEvent): Promise<any> {
  console.log('Event:', JSON.stringify(event, null, 2));

  const { data, success, error } = commitFormSchema.safeParse(JSON.parse(event.body || '{}'));

  if (!success) {
    console.error('Validation error:', data);

    const treeError = z.treeifyError(error);

    return {
      statusCode: 400,
      body: JSON.stringify({ 
        message: `Invalid input data: ${error.message}`,
        type: error.type,
        errors: treeError.errors
      })
    };
  }

  const { user, solution, user_access_token } = data;

  try {
    const octokit = new Octokit({auth: user_access_token});


    const mainBranchRef = await octokit.git.getRef({
      owner: user.owner,
      repo: user.repo_name,
      ref: 'heads/main',
    });

    console.log('get ref passed');

    const commitSha = mainBranchRef.data.object.sha;

    const latestCommit = await octokit.git.getCommit({
      owner: user.owner,
      repo: user.repo_name,
      commit_sha: commitSha,
    });

    console.log('get commit passed');

    // Commit tree setup
    const basePath = solution.problem_id ? solution.problem_name : `${solution.problem_id}-${solution.problem_name}`;

    const languageExtension = mapProblemLanguageToExtension(solution.language);

    let description: string = "";
    
    if (solution.problem_link) {
      description += `<h3><a href=${solution.problem_link}>${solution.problem_name}</a></h3>\n`;
    }

    if (solution.description) {
      description += `<p>${solution.description}</p>\n`;
    }

    if (description === "") {
      description += solution.description;
    }

    const tree = await octokit.git.createTree({
      owner: user.owner,
      repo: user.repo_name,
      tree: [
        {
          path: `${basePath}/${solution.problem_name}.${languageExtension}`,
          type: 'blob',
          content: solution.code,
          mode: '100755',
        },
        {
          path: `${basePath}/README.md`,
          type: 'blob',
          content: description,
          mode: '100644'
        },
        {
          path: `${basePath}/NOTES.md`,
          type: 'blob',
          content: solution.notes ?? "",
          mode: '100644'
        }
      ],
      base_tree: latestCommit.data.tree.sha,
    });

    console.log('create tree passed');

    const newCommitResponse = await octokit.git.createCommit({
      owner: user.owner,
      repo: user.repo_name,
      message: `ADD ${solution.problem_name}`,
      tree: tree.data.sha,
      parents: [latestCommit.data.sha],
    });

    console.log('create commit passed');


    await octokit.git.updateRef({
      owner: user.owner,
      repo: user.repo_name,
      ref: 'heads/main',
      sha: newCommitResponse.data.sha,
    });


    return {
      statusCode: 200,
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'OPTIONS,POST,GET',
        'Access-Control-Allow-Headers': '*',
      },
      body: JSON.stringify({
        message: `Commit processed successfully for problem: ${solution.problem_name}`,
        user,
        solution,
      }),
    };
  } catch (error) {
    console.error('Error commiting problem:', error);
    
    return {
      statusCode: 500,
      body: JSON.stringify({ 
        message: `Failed to commit solution for problem ${solution.problem_name}. Reason: ${error instanceof Error ? error.message : 'Unknown error'}`, 
      })
    };
  }
}
