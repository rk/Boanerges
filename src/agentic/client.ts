import { call, event } from "mygo-runtime";

export interface OllamaHealth {
  ok: boolean;
  error?: string;
  models?: string[];
}

export interface ProjectManifest {
  kind: "dictionary" | "topical";
  model: string;
}

export interface RunStep {
  at: string;
  role: string;
  summary: string;
}

export interface RunState {
  status: "idle" | "running" | "paused" | "failed";
  job: string;
  instruction: string;
  params?: unknown;
  steps: RunStep[];
  error?: string;
  updatedAt: string;
}

export interface ProjectSummary {
  folder: string;
  manifest: ProjectManifest;
  run: RunState;
  sourceCount: number;
  entryCount: number;
  topicCount: number;
}

export interface EntrySummary {
  headword: string;
  file: string;
}

export interface TopicSummary {
  slug: string;
  title: string;
}

export const agenticEvents = {
  runStateChanged: event<{}>("agenticRunStateChanged"),
};

export class Agentic {
  pingOllama(): Promise<OllamaHealth> {
    return call("Agentic.PingOllama");
  }
  pickProjectFolder(): Promise<string> {
    return call("Agentic.PickProjectFolder");
  }
  initProject(folder: string, kind: string): Promise<void> {
    return call("Agentic.InitProject", folder, kind);
  }
  loadProjectSummary(folder: string): Promise<ProjectSummary> {
    return call("Agentic.LoadProjectSummary", folder);
  }
  listSources(folder: string): Promise<string[]> {
    return call("Agentic.ListSources", folder);
  }
  listEntries(folder: string): Promise<EntrySummary[]> {
    return call("Agentic.ListEntries", folder);
  }
  listTopics(folder: string): Promise<TopicSummary[]> {
    return call("Agentic.ListTopics", folder);
  }
  getRunState(folder: string): Promise<RunState> {
    return call("Agentic.GetRunState", folder);
  }
  setProjectModel(folder: string, model: string): Promise<void> {
    return call("Agentic.SetProjectModel", folder, model);
  }
  startDictionaryRun(folder: string, instruction: string): Promise<void> {
    return call("Agentic.StartDictionaryRun", folder, instruction);
  }
  analyzeSourceFormats(folder: string): Promise<void> {
    return call("Agentic.AnalyzeSourceFormats", folder);
  }
  startTopicalFromReference(
    folder: string,
    bookId: string,
    chapter: number,
    verse: number,
    instruction: string,
  ): Promise<void> {
    return call("Agentic.StartTopicalFromReference", folder, bookId, chapter, verse, instruction);
  }
  startTopicalFromTopic(folder: string, topic: string, instruction: string): Promise<void> {
    return call("Agentic.StartTopicalFromTopic", folder, topic, instruction);
  }
  pauseRun(folder: string): Promise<void> {
    return call("Agentic.PauseRun", folder);
  }
  resumeRun(folder: string): Promise<void> {
    return call("Agentic.ResumeRun", folder);
  }
  readRawRunLog(folder: string, maxBytes: number): Promise<string> {
    return call("Agentic.ReadRawRunLog", folder, maxBytes);
  }
}

export const agentic = new Agentic();
