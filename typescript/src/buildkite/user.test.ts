import { registerUserTools } from './user';

describe('registerUserTools', () => {
  it('registers current_user tool on the server', () => {
    const tools: any[] = [];
    const mockServer = {
      tool: jest.fn((name, desc, schema, handler) => {
        tools.push({ name, desc, schema, handler });
      })
    };
    const mockClient = {};
    registerUserTools(mockServer as any, mockClient as any);
    const toolNames = tools.map(t => t.name);
    expect(toolNames).toContain('current_user');
    expect(tools.length).toBe(1);
  });
});

