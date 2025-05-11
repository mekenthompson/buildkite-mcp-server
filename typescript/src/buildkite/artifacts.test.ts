import { registerArtifactTools } from './artifacts';

describe('registerArtifactTools', () => {
  it('registers list_artifacts and get_artifact tools on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerArtifactTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('list_artifacts');
    expect(toolNames).toContain('get_artifact');
    expect(tools.length).toBe(2);
  });
});
