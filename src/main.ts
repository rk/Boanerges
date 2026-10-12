import { mount } from "svelte";
import App from "./App.svelte";
import "./style.css";

async function bootstrap(): Promise<void> {
  if (import.meta.env.VITE_E2E === "1") {
    const { installMygoMock } = await import("./e2e/mock-mygo");
    installMygoMock();
  }

  mount(App, { target: document.getElementById("app")! });
}

void bootstrap();
