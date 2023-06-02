import { Injectable, InternalServerErrorException } from '@nestjs/common';
import axios from 'axios';
import { get } from 'src/config';
import { SecretsService } from 'src/secrets/secrets.service';
import { Organization } from 'src/typeorm';

@Injectable()
export class GithubApiService {
  private readonly baseUri: string = 'https://api.github.com/repos';
  private readonly headers = async (org: Organization) => ({
    Accept: 'application/vnd.github+json',
    Authorization: `Bearer ${await this.getGITPAT(org)}`,
    'X-GitHub-Api-Version': '2022-11-28',
  });

  constructor(private readonly secretSvc: SecretsService){}

  private getURL(owner: string, repo: string, filePath: string) {
    return `${this.baseUri}/${owner}/${repo}/contents/${filePath}`;
  }

  private async getFileSHA(org: Organization, owner: string, repo: string, filePath: string) {
    const url = this.getURL(owner, repo, filePath);
    const { data } = await axios.get<{ sha: string; content: string }>(url, {
      headers: await this.headers(org),
    });
    return data;
  }

  public setAlternateFlag(yamlString: string) {
    const decodedYaml = Buffer.from(yamlString, 'base64').toString();
    console.log(decodedYaml);
    if (decodedYaml.includes('teardown: true')) {
      return Buffer.from(
        decodedYaml.replace('teardown: true', 'teardown: false')
      ).toString('base64');
    }
    return Buffer.from(
      decodedYaml.replace('teardown: false', 'teardown: true')
    ).toString('base64');
  }

  public async gitCommit(org: Organization, owner: string, repo: string, filePath: string) {
    const { sha, content } = await this.getFileSHA(org, owner, repo, filePath);
    const payload = {
      message: 'api testing for playground ',
      committer: { name: 'playground', email: 'playground@cloudknit.io' },
      content: this.setAlternateFlag(content),
      sha,
    };

    try {
      const { data } = await axios.put(
        this.getURL(owner, repo, filePath),
        payload,
        {
          headers: await this.headers(org),
        }
      );

      const { commit } = data;
      const { html_url } = commit;
      return {
        status: 'success',
        html_url,
      };
    } catch (err) {
      throw new InternalServerErrorException(
        'There was an error while pushing the commit to git'
      );
    }
  }

  private async getGITPAT(org: Organization) {
    return await this.secretSvc.getSsmSecret(org, 'playground-git-token');
  }
}
