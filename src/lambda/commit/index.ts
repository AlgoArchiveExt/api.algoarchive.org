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
    accessToken: z.string()
})

export const handler = async function (event: APIGatewayProxyEvent): Promise<any> {
  console.log('Event:', JSON.stringify(event, null, 2));

  const { data, success, error } = commitFormSchema.safeParse(JSON.parse(event.body || '{}'));

  if (!success) {
    console.error('Validation error:', data);
    return {
      statusCode: 400,
      body: JSON.stringify({ 
        message: `Invalid input data: ${error.message}`, 
        errors: z.treeifyError(error).errors
      })
    };
  }

  const { user, solution, accessToken } = data;

  try {


    const octokit = new Octokit({auth: accessToken});


    const mainBranchRef = await octokit.git.getRef({
      owner: user.owner,
      repo: user.repoName,
      ref: 'heads/main',
    });

    const commitSha = mainBranchRef.data.object.sha;

    const latestCommit = await octokit.git.getCommit({
      owner: user.owner,
      repo: user.repoName,
      commit_sha: commitSha,
    });

    // Commit tree setup
    const basePath = solution.problemId ? solution.problemName : `${solution.problemId}-${solution.problemName}`;

    const languageExtension = mapProblemLanguageToExtension(solution.language);

    let description: string = ``;
    
    if (solution.problemLink) {
      description += `<h3><a href=${solution.problemLink}>${solution.problemName}</a></h3>\n`;
    }

    if (solution.description) {
      description += `<p>${solution.description}</p>\n`;
    }

    description += solution.description;

    const tree = await octokit.git.createTree({
      owner: user.owner,
      repo: user.repoName,
      tree: [
        {
          path: `${basePath}/${solution.problemName}.${languageExtension}`,
          type: 'blob',
          content: solution.code,
          mode: '100755'
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
          content: solution.notes,
          mode: '100644'
        }
      ],
      base_tree: latestCommit.data.tree.sha,
    });

    const newCommitResponse = await octokit.git.createCommit({
      owner: user.owner,
      repo: user.repoName,
      message: `ADD ${solution.problemName}`,
      tree: tree.data.sha,
      parents: [latestCommit.data.sha],
    });


    await octokit.git.updateRef({
      owner: user.owner,
      repo: user.repoName,
      ref: 'heads/main',
      sha: newCommitResponse.data.sha,
    });


    return {
      statusCode: 200,
      body: JSON.stringify({
        message: `Commit processed successfully for problem: ${solution.problemName}`,
        user,
        solution,
      }),
    };
  } catch (error) {
    console.error('Error commiting problem:', error);
    
    return {
      statusCode: 500,
      body: JSON.stringify({ 
        message: `Failed to commit solution for problem ${solution.problemName}. Reason: ${error instanceof Error ? error.message : 'Unknown error'}`, 
      })
    };
  }
}
