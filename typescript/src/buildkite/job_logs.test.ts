import { registerJobTools } from './job_logs';

describe('registerJobTools', () => {
  it('registers get_job_logs tool on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerJobTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('get_job_logs');
    expect(tools.length).toBe(1);
  });
});
