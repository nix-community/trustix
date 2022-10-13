import type { Component } from "solid-js";
import { lazy, onMount } from "solid-js";
import { Routes, Route, A } from "@solidjs/router";
import { themeChange } from "theme-change";

const Derivation = lazy(() => import("./pages/drv"));

const App: Component = () => {
  onMount(async () => {
    themeChange();
  });

  return (
    <>
      <div className="bg-base-200">
        <nav id="main-nav">
          <A href="/channels">Channels</A>
          <A href="/channel">Channels</A>
          <A href="/drv">Derivation</A>
        </nav>

        <Routes>
          <Route path="/drv" component={Derivation} />
        </Routes>
      </div>
    </>
  );
};

export default App;
