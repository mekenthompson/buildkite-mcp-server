import fetch, { Response, RequestInit, HeadersInit } from 'node-fetch';

// Base URL for Buildkite API
const API_BASE_URL = process.env.BUILDKITE_API_BASE_URL || 'https://api.buildkite.com/v2';

// Type definitions for Buildkite API responses
export interface BuildkitePipeline {
  id: string;
  url: string;
  web_url: string;
  name: string;
  slug: string;
  repository: string;
  branch_configuration: string | null;
  default_branch: string;
  provider: {
    id: string;
    webhook_url: string;
    settings: Record<string, any>;
  };
  skip_queued_branch_builds: boolean;
  skip_queued_branch_builds_filter: string | null;
  cancel_running_branch_builds: boolean;
  cancel_running_branch_builds_filter: string | null;
  builds_url: string;
  badge_url: string;
  created_at: string;
  scheduled_builds_count: number;
  running_builds_count: number;
  scheduled_jobs_count: number;
  running_jobs_count: number;
  waiting_jobs_count: number;
  visibility: string;
  [key: string]: any;
}

export interface BuildkiteUser {
  id: string;
  name: string;
  email: string;
  avatar_url: string;
  created_at: string;
  [key: string]: any;
}

export interface BuildkiteBuild {
  id: string;
  url: string;
  web_url: string;
  number: number;
  state: string;
  blocked: boolean;
  message: string;
  commit: string;
  branch: string;
  tag: string | null;
  created_at: string;
  scheduled_at: string;
  started_at: string | null;
  finished_at: string | null;
  [key: string]: any;
}

export interface BuildkiteArtifact {
  id: string;
  job_id: string;
  url: string;
  download_url: string;
  state: string;
  path: string;
  dirname: string;
  filename: string;
  mime_type: string;
  file_size: number;
  sha1sum: string;
  [key: string]: any;
}

export interface BuildkiteOrganization {
  id: string;
  url: string;
  web_url: string;
  name: string;
  slug: string;
  [key: string]: any;
}

// Client class for making authenticated requests to Buildkite API
export class BuildkiteClient {
  private apiToken: string;
  private userAgent: string;

  constructor(apiToken: string, version: string) {
    this.apiToken = apiToken;
    this.userAgent = `buildkite-mcp-server/${version}`;
  }

  // Helper method to build request headers
  private getHeaders(): HeadersInit {
    return {
      'Authorization': `Bearer ${this.apiToken}`,
      'User-Agent': this.userAgent,
      'Content-Type': 'application/json',
    };
  }

  async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${path}`;
    const headers: HeadersInit = { ...this.getHeaders(), ...options.headers };
    
    try {
      const response = await fetch(url, { ...options, headers });
      
      if (!response.ok) {
        const errorBody = await response.text();
        throw new Error(`Buildkite API error (${response.status}): ${errorBody}`);
      }
      
      // Return empty object if no content
      if (response.status === 204) {
        return {} as T;
      }
      
      return await response.json() as T;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error(`Unknown error communicating with Buildkite API: ${String(error)}`);
    }
  }

  // Pipeline-related methods
  async getPipeline(org: string, pipelineSlug: string): Promise<BuildkitePipeline> {
    return this.request<BuildkitePipeline>(`/organizations/${org}/pipelines/${pipelineSlug}`);
  }

  async listPipelines(org: string, page = 1, perPage = 30): Promise<BuildkitePipeline[]> {
    return this.request<BuildkitePipeline[]>(`/organizations/${org}/pipelines?page=${page}&per_page=${perPage}`);
  }

  // Build-related methods
  async listBuilds(org: string, pipelineSlug: string, page = 1, perPage = 30): Promise<BuildkiteBuild[]> {
    return this.request<BuildkiteBuild[]>(`/organizations/${org}/pipelines/${pipelineSlug}/builds?page=${page}&per_page=${perPage}`);
  }

  async getBuild(org: string, pipelineSlug: string, buildNumber: string): Promise<BuildkiteBuild> {
    return this.request<BuildkiteBuild>(`/organizations/${org}/pipelines/${pipelineSlug}/builds/${buildNumber}`);
  }

  // Job-related methods
  async getJobLogs(org: string, pipelineSlug: string, buildNumber: string, jobId: string): Promise<any> {
    return this.request<any>(`/organizations/${org}/pipelines/${pipelineSlug}/builds/${buildNumber}/jobs/${jobId}/log`);
  }

  // Artifact-related methods
  async listArtifacts(org: string, pipelineSlug: string, buildNumber: string, jobId: string, page = 1, perPage = 30): Promise<BuildkiteArtifact[]> {
    return this.request<BuildkiteArtifact[]>(
      `/organizations/${org}/pipelines/${pipelineSlug}/builds/${buildNumber}/jobs/${jobId}/artifacts?page=${page}&per_page=${perPage}`
    );
  }

  async getArtifact(org: string, pipelineSlug: string, buildNumber: string, jobId: string, artifactId: string): Promise<BuildkiteArtifact> {
    return this.request<BuildkiteArtifact>(
      `/organizations/${org}/pipelines/${pipelineSlug}/builds/${buildNumber}/jobs/${jobId}/artifacts/${artifactId}`
    );
  }

  async downloadArtifact(url: string): Promise<Response> {
    return fetch(url, { headers: this.getHeaders() });
  }

  // User-related methods
  async getCurrentUser(): Promise<BuildkiteUser> {
    return this.request<BuildkiteUser>('/user');
  }

  // Organization-related methods
  async listOrganizations(page = 1, perPage = 30): Promise<BuildkiteOrganization[]> {
    return this.request<BuildkiteOrganization[]>(`/organizations?page=${page}&per_page=${perPage}`);
  }

  async getUserTokenOrganization(): Promise<BuildkiteOrganization> {
    const orgs = await this.listOrganizations();
    if (Array.isArray(orgs) && orgs.length > 0) {
      return orgs[0];
    }
    throw new Error('No organizations found for this token');
  }

  // Access token methods
  async getAccessToken(): Promise<any> {
    return this.request<any>('/access-token');
  }
}

// Factory function to create a BuildkiteClient instance
export function createBuildkiteClient(apiToken: string, version: string): BuildkiteClient {
  return new BuildkiteClient(apiToken, version);
}