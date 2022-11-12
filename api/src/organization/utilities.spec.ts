import { getGithubOrgFromRepoUrl } from "./utilities";

describe('Organization Service', () => {
  const badRepoException = 'GitHub repo url is not formatted properly';
  const orgName = 'some-random-org';

  beforeEach(async () => { });

  describe('getGithubOrgFromRepoUrl ssh', () => {
    const repoUrl = `git@github.com:${orgName}/hello-world.git`;

    it('should return the org from the url', async () => {
      const org1 = getGithubOrgFromRepoUrl(repoUrl);

      expect(org1).toStrictEqual(orgName);
    });

    it('should throw exception', async () => {
      expect(getGithubOrgFromRepoUrl).toThrow();
    });

    it('should throw exception', async () => {
      expect(() => {
        getGithubOrgFromRepoUrl('git@github.comsome-random-org/hello-world.git');
      }
      ).toThrow(badRepoException);
    });

    it('should throw exception', async () => {
      expect(() => {
        getGithubOrgFromRepoUrl('git@github.com:some-random-org/hello-world');
      }
      ).toThrow(badRepoException);
    });
  });

  describe('getGithubOrgFromRepoUrl https', () => {
    const repoUrl = `https://github.com/${orgName}/hello-world.git`;

    it('should return the org from the url', async () => {
      const org1 = getGithubOrgFromRepoUrl(repoUrl);

      expect(org1).toStrictEqual(orgName);
    });

    it('should throw exception', async () => {
      expect(() => {
        getGithubOrgFromRepoUrl('https://github.com/some-random-orghello-world.git');
      }
      ).toThrow(badRepoException);
    });
  });
});
