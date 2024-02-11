import { createContext } from "react";

export const PreviewRequestContext = createContext<((fileUrl: string, fileMimetype:string) => void) | null>(null);
