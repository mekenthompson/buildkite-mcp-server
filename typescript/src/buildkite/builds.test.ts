import { registerBuildTools } from './builds';

describe('registerBuildTools', () => {
  it('registers list_builds and get_build tools on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerBuildTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('list_builds');
    expect(toolNames).toContain('get_build');
    expect(tools.length).toBe(2);
  });
});