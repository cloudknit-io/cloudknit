export function getGithubOrgFromRepoUrl(repoUrl: string) : string {
  if (!repoUrl || !repoUrl.endsWith('.git')) {
    return null;
  }
  
  if (repoUrl.startsWith('git@')) {
    // git@github.com:some-random-org/hello-world.git
    const allParts = repoUrl.split(':');

    if (allParts.length !== 2) {
      return null;
    }

    const orgAndRepo = allParts[1].split('/');

    if (orgAndRepo.length !== 2) {
      return null;
    }

    return orgAndRepo[0];
  } else if (repoUrl.startsWith('https://')) {
    // https://github.com/some-random-org/hello-world.git
    const allParts = repoUrl.split('/');

    if (allParts.length !== 5) {
      return null;
    }

    return allParts[3];
  } else {
    return null;
  }    
};
