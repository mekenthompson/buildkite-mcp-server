import { registerOrganizationTools } from './organizations';

describe('registerOrganizationTools', () => {
  it('registers user_token_organization and user_token_organization_prompt tools on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      }),
      prompt: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerOrganizationTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('user_token_organization');
    expect(toolNames).toContain('user_token_organization_prompt');
    expect(tools.length).toBe(2);
  });
});