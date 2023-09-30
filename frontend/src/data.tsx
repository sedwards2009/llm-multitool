
export interface SessionOverview {
  sessionSummaries: SessionSummary[];
}

export interface SessionSummary {
  id: string;
  title: string;
  creationTimestamp: string;
}

export interface Root {
  sessions: Session[];
}

export interface Session {
  id: string;
  creationTimestamp: string;
  title: string;
  prompt: string;
  responses: Response[];
  modelSettings: ModelSettings;
}

export interface ModelSettings {
  modelId: string;
  templateId: string;
  presetId: string;
}

export type ResponseStatus = "Done" | "Pending" | "Running" | "Error";

export interface Response {
  id: string;
  creationTimestamp: string;
  status: ResponseStatus;
  messages: Message[];
  modelSettingsSnapshot: ModelSettingsSnapshot;
}

export interface ModelSettingsSnapshot {
	modelName: string;
	templateName: string;
	presetName: string;
}

export type Role = "User" | "Assistant";

export interface Message {
  id: string;
  role: Role;
  text: string;
}

export interface Model {
  id: string;
  name: string;
}

export interface ModelOverview {
  models: Model[];
}

export interface Template {
  id: string;
  name: string;
}

export interface TemplateOverview {
  templates: Template[];
}

export interface Preset {
  id: string;
  name: string;
}

export interface PresetOverview {
  presets: Preset[];
}
