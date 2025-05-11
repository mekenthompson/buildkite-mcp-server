import { registerAccessTokenTools } from './access_token';

describe('registerAccessTokenTools', () => {
  it('registers access_token tool on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerAccessTokenTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('access_token');
    expect(tools.length).toBe(1);
  });
});