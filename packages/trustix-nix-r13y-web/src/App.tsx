import type { Component } from "solid-js";
import { lazy, onMount, For } from "solid-js";
import { Routes, Route, A, Navigate } from "@solidjs/router";
import { themeChange } from "theme-change";

const Derivation = lazy(() => import("./pages/drv"));
const Attrs = lazy(() => import("./pages/attrs"));

const links: {
  title: string;
  href: string;
}[] = [
  {
    title: "Attributes",
    href: "/attrs",
  },
];

const App: Component = () => {
  onMount(async () => {
    themeChange();
  });

  return (
    <>
      <div className="bg-base-200 min-h-screen">
        <nav id="main-nav" className="bg-primary p-2">
          <For each={links}>
            {(l) => (
              <A className="text-primary-content m-2" href={l.href}>
                {l.title}
              </A>
            )}
          </For>

          <button
            className="float-right"
            data-toggle-theme="dark,light"
            data-act-class="ACTIVECLASS"
          >
            <span className="" data-tip="toggle dark mode">
              💡
            </span>
          </button>
        </nav>

        <div>
          <Routes>
            <Route path="/" component={<Navigate href="/attrs" />} />
            <Route path="/attrs" component={Attrs} />
            <Route path="/drv" component={Derivation} />
          </Routes>
        </div>
      </div>
    </>
  );
};

export default App;
