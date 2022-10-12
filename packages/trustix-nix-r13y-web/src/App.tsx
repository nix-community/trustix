import type { Component } from 'solid-js';
import { lazy } from 'solid-js';
import { Routes, Route, A } from "@solidjs/router"

const Derivation = lazy(() => import("./pages/drv"))

const App: Component = () => {
  return (
    <>
      <p>Trustix r13y</p>

      <nav>
        <A href="/channels">Channels</A>
        <A href="/channel">Channels</A>
        <A href="/drv">Derivation</A>
      </nav>

      <Routes>
        <Route path="/drv" component={Derivation} />
      </Routes>

    </>
  );
};

export default App;
