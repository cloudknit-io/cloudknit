import { generateParams } from './api';

describe('Generate Params', () => {
  beforeEach(async () => {});

  it('should return formatted string', () => {
    const obj = {
      orgName: 'brad-org-2'
    };

    const resp = generateParams(obj);

    expect(resp).toBeDefined();
    expect(resp).toHaveLength(1);
    expect(resp[0]).toStrictEqual('orgName=brad-org-2');
  });

  it('should return formatted string', () => {
    const obj = {
      orgName: 'brad-org-2',
      orgId: 1,
      orgMeta: 'brad-meta'
    };

    const resp = generateParams(obj);

    expect(resp).toBeDefined();
    expect(resp).toHaveLength(3);
    expect(resp).toContain('orgName=brad-org-2');
    expect(resp).toContain('orgId=1');
    expect(resp).toContain('orgMeta=brad-meta');
  });
});
