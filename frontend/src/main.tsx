import * as React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { DomainProvider } from "./context/DomainContext";
import App from "./App";
import "./index.css"; // or "./main.css"

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <DomainProvider>
        <App />
      </DomainProvider>
    </BrowserRouter>
  </React.StrictMode>
);
