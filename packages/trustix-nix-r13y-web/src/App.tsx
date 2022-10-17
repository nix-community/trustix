// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import type { Component } from "solid-js";
import { lazy, onMount, For } from "solid-js";
import { Routes, Route, A } from "@solidjs/router";
import { themeChange } from "theme-change";

const Derivation = lazy(() => import("./pages/drv"));
const Attrs = lazy(() => import("./pages/attrs"));
const Diff = lazy(() => import("./pages/diff"));

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
      <div class="bg-base-200 min-h-screen">
        <nav id="main-nav" class="bg-primary p-2">
          <For each={links}>
            {(l) => (
              <A class="text-primary-content m-2" href={l.href}>
                {l.title}
              </A>
            )}
          </For>

          <button
            class="float-right"
            data-toggle-theme="dark,light"
            data-act-class="ACTIVECLASS"
          >
            <span class="" data-tip="toggle dark mode">
              💡
            </span>
          </button>
        </nav>

        <div class="flex justify-evenly place-items-center">
          <Routes>
            <Route path="/attrs" component={Attrs} />
            <Route path="/drv" component={Derivation} />
            <Route path="/diff" component={Diff} />
            {/* <Route path="/" component={<Navigate href="/attrs" />} /> */}
          </Routes>
        </div>
      </div>
    </>
  );
};

export default App;
