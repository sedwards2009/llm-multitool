import * as _ from "lodash";
import { ModelOverview, ModelSettings, PresetOverview, Session, SessionOverview, TemplateOverview } from "./data";

export interface LoaderResult {
  sessionOverview: SessionOverview;
  session?: Session;
}

const windowProtocol = window.location.protocol;
const SERVER_BASE_URL = `${window.location.protocol}//${window.location.hostname}:${window.location.port}/api`;
// const SERVER_BASE_URL = "http://localhost:8080/api";
const WEBSOCKET_SERVER_BASE_URL = `${windowProtocol === "http:" ? "ws" : "wss"}://${window.location.hostname}:${window.location.port}/api`;
// const WEBSOCKET_SERVER_BASE_URL = "ws://localhost:8080/api";

export async function loadModelOverview(): Promise<ModelOverview> {
  const response = await fetch(`${SERVER_BASE_URL}/model`);
  try {
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return {
      models: []
    };
  }
}

export async function scanModels(): Promise<ModelOverview> {
  const response = await fetch(`${SERVER_BASE_URL}/model/scan`, {
    method: "POST"
  });
  try {
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return {
      models: []
    };
  }
}

export async function loadSessionOverview(): Promise<SessionOverview> {
  const response = await fetch(`${SERVER_BASE_URL}/session`);
  try {
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return {
      sessionSummaries: []
    };
  }
}

export async function loadPresetOverview(): Promise<PresetOverview> {
  const response = await fetch(`${SERVER_BASE_URL}/preset`);
  try {
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return {
      presets: []
    };
  }
}

export async function loadTemplateOverview(): Promise<TemplateOverview> {
  const response = await fetch(`${SERVER_BASE_URL}/template`);
  try {
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return {
      templates: []
    };
  }
}

export async function loadSession(sessionId: string): Promise<Session | null> {
  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}`);
  try {
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
  return null;
}

export async function newSession(defaults?: ModelSettings | null): Promise<Session | null> {
  const fetchOptions: RequestInit = {
    method: "POST",
  };
  if (defaults != null) {
    fetchOptions.headers = {"Content-Type": "application/json"};
    fetchOptions.body = JSON.stringify(defaults);
  }

  const response = await fetch(`${SERVER_BASE_URL}/session`, fetchOptions);
  try {
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
  return null;
}

export async function deleteSession(sessionId: string): Promise<void> {
  await fetch(`${SERVER_BASE_URL}/session/${sessionId}`, {method: "DELETE"});
}

async function putSessionProperty(sessionId: string, propertyName: string, value: any): Promise<void> {
  await fetch(`${SERVER_BASE_URL}/session/${sessionId}/${propertyName}`,
    {
      method: "PUT",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(value)
    });
}

// This is all about throttling the REST calls done to update the prompt on the
// server when the user is typing a prompt. We don't need to write every keypress.
let saveSessionPromptQueue = new Map<string, string>();

export function setSessionPrompt(session: Session, prompt: string): Session {
  saveSessionPromptQueue.set(session.id, prompt);
  flushSessionPromptQueue();
  return {...session, prompt};
}

async function processSessionPromptQueue(): Promise<void> {
  const workingSessionPromptQueue = saveSessionPromptQueue;
  saveSessionPromptQueue = new Map<string, string>();
  for (const [sessionId, prompt] of workingSessionPromptQueue.entries()) {
    await putSessionProperty(sessionId, "prompt", {value: prompt});
  }
}

const flushSessionPromptQueue = _.throttle(() => {
  (async () => {
    processSessionPromptQueue();
  })();
}, 1000);

async function flushQueues(): Promise<void> {
  await processSessionPromptQueue();
}

export function setSessionModel(session: Session, modelId: string): Session {
  const newModelSettings: ModelSettings = {...session.modelSettings, modelId };
  putSessionProperty(session.id, "modelSettings", newModelSettings);
  return {...session, modelSettings: newModelSettings};
}

export function setSessionTemplate(session: Session, templateId: string): Session {
  const newModelSettings: ModelSettings = {...session.modelSettings, templateId: templateId };
  putSessionProperty(session.id, "modelSettings", newModelSettings);
  return {...session, modelSettings: newModelSettings};
}

export function setSessionPreset(session: Session, presetId: string): Session {
  const newModelSettings: ModelSettings = {...session.modelSettings, presetId: presetId };
  putSessionProperty(session.id, "modelSettings", newModelSettings);
  return {...session, modelSettings: newModelSettings};
}

export async function uploadFileToSession(session: Session, file: File): Promise<void> {
  const formData = new FormData();
  formData.append("file", file);

  if (file.type != null) {
    formData.append("mimeType", file.type);
  }
  if (file.name) {
    formData.append("filename", file.name);
  }

  await fetch(`${SERVER_BASE_URL}/session/${session.id}/file`,
    {
      method: "POST",
      body: formData,
    });
}

export function fileURL(sessionId: string, fileId: string): string {
  return `${SERVER_BASE_URL}/session/${sessionId}/file/${fileId}`;
}

export async function deleteSessionAttachedFile(sessionId: string, filename: string): Promise<boolean> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}/file/${filename}`, {method: "DELETE"});
  return response.ok;
}

export async function newResponse(session: Session): Promise<Response | null> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${session.id}/response`, {method: "POST"});
  try {
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
  return null;
}

export async function abortResponse(sessionId: string, responseId: string): Promise<boolean> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}/response/${responseId}/abort`, {method: "POST"});
  return response.ok;
}

export async function deleteResponse(sessionId: string, responseId: string): Promise<boolean> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}/response/${responseId}`, {method: "DELETE"});
  return response.ok;
}

export async function deleteResponseMessage(sessionId: string, responseId: string, messageId: string): Promise<boolean> {
  await flushQueues();
  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}/response/${responseId}/message/${messageId}`,
    {method: "DELETE"});
  return response.ok;
}

export async function newMessage(session: Session, responseId: string, reply: string): Promise<void> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${session.id}/response/${responseId}/message`,
    {
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({value: reply}),
      method: "POST"
    });
  try {
    if (response.ok) {
      return;
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
}

export async function continueMessage(session: Session, responseId: string): Promise<void> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${session.id}/response/${responseId}/continue`,
    {
      headers: {"Content-Type": "application/json"},
      method: "POST"
    });
  try {
    if (response.ok) {
      return;
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
}

export enum SessionMonitorState {
  IDLE,
  CONNECTING,
  CONNECTED,
  WAITING_TO_RECONNECT,
}

const DEFAULT_RECONNECT_DELAY_MS = 100;

export class SessionMonitor {
  #sessionId = "";
  #socket: WebSocket | null = null;
  #state = SessionMonitorState.IDLE;
  #onChange: ((message: string) => void)  | null = null;
  #onStateChange: ((statue: SessionMonitorState) => void)  | null = null;
  #reconnectDelayMs = DEFAULT_RECONNECT_DELAY_MS;

  constructor(sessionId: string, onChange: (message: string) => void,
      onStateChange: (state: SessionMonitorState) => void) {
    this.#sessionId = sessionId;
    this.#onChange = onChange;
    this.#onStateChange = onStateChange;
  }

  state(): SessionMonitorState {
    return this.#state;
  }

  #setState(state: SessionMonitorState): void {
    this.#state = state;
    if (this.#onStateChange != null) {
      console.log(`SessionMonitor ${state}`);
      this.#onStateChange(state);
    }
  }

  start(): void {
    this.#connect();
  }

  #connect(): void {
    console.log(`Connecting`);
    this.#setState(SessionMonitorState.CONNECTING);
    this.#socket = new WebSocket(`${WEBSOCKET_SERVER_BASE_URL}/session/${this.#sessionId}/changes`)
    this.#socket.addEventListener("message", (event) => {
      if (this.#onChange != null) {
        this.#onChange(event.data);
      }
    });
    this.#socket.addEventListener("open", () => {
      this.#setState(SessionMonitorState.CONNECTED);
      this.#reconnectDelayMs = DEFAULT_RECONNECT_DELAY_MS;
    });
    this.#socket.addEventListener("close", () => {
      if (this.#state === SessionMonitorState.IDLE) {
        return;
      }
      this.#reconnect();
    });
    this.#socket.addEventListener("error", (e) => {
      console.log(`Websocket error for sessionId ${this.#sessionId}`, e);
      if (this.#socket !== null) {
        this.#socket.close();
      }
    });
  }

  #reconnect(): void {
    this.#setState(SessionMonitorState.WAITING_TO_RECONNECT);
    console.log(`Reconnecting ${this.#reconnectDelayMs}ms`);
    setTimeout(() => {
      if (this.#state === SessionMonitorState.IDLE) {
        return;
      }
      this.#connect();
    }, this.#reconnectDelayMs);
    this.#reconnectDelayMs = Math.min(this.#reconnectDelayMs * 2, 5000);
  }

  stop(): void {
    this.#state = SessionMonitorState.IDLE;
    if (this.#socket == null) {
      return;
    }
    this.#socket.close();
    this.#socket = null;
  }
}