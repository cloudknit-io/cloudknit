import { Injectable, InternalServerErrorException } from '@nestjs/common';
import axios from 'axios';
import { get } from 'src/config';

@Injectable()
export class GithubApiService {
  private readonly baseUri: string = 'https://api.github.com/repos';
  private readonly headers = {
    Accept: 'application/vnd.github+json',
    Authorization: `Bearer ${Buffer.from(get().github.personalAccessToken, 'base64').toString()}`,
    'X-GitHub-Api-Version': '2022-11-28',
  };

  private getURL(owner: string, repo: string, filePath: string) {
    return `${this.baseUri}/${owner}/${repo}/contents/${filePath}`;
  }

  private async getFileSHA(owner: string, repo: string, filePath: string) {
    const url = this.getURL(owner, repo, filePath);
    const { data } = await axios.get<{ sha: string; content: string }>(url, {
      headers: this.headers,
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

  public async gitCommit(owner: string, repo: string, filePath: string) {
    const { sha, content } = await this.getFileSHA(owner, repo, filePath);
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
          headers: this.headers,
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
}
