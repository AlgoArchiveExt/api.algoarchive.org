export function mapProblemLanguageToExtension(language: string): string {
  const extensions: Record<string, string> = {
    'python': 'py',
    'python3': 'py',
    'pythondata': 'py',
    'javascript': 'js',
    'typescript': 'ts',
    'mysql': 'sql',
    'postgresql': 'sql',
    'oraclesql': 'sql',
    'mssql': 'sql',
    'java': 'java',
    'c': 'c',
    'cpp': 'cpp',
    'csharp': 'cs',
    'go': 'go',
    'ruby': 'rb',
    'rust': 'rs',
    'kotlin': 'kt',
    'swift': 'swift',
    'php': 'php',
    'scala': 'scala',
    'dart': 'dart',
    'elixir': 'ex',
    'racket': 'rkt',
    'erlang': 'erl',

  };

  return extensions[language.toLowerCase()] || 'txt';
}