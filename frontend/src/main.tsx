import React from 'react';
import ReactDOM from 'react-dom/client';
import { App } from './app.tsx';
import './theme/main.scss';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import { SessionLoader, SessionOverviewLoader } from './dataloading.tsx';



const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    loader: SessionOverviewLoader
  },
  {
    path: "session/:sessionId",
    element: <App />,
    loader: SessionLoader
  }
]);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
