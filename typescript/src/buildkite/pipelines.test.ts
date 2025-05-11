import { registerPipelineTools } from './pipelines';

describe('registerPipelineTools', () => {
  it('registers list_pipelines and get_pipeline tools on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerPipelineTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('list_pipelines');
    expect(toolNames).toContain('get_pipeline');
    expect(tools.length).toBe(2);
  });
});