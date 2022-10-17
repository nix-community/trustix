// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import { Component, createResource, Show, Suspense } from "solid-js";
import { useSearchParams } from "@solidjs/router";

import { DiffResponse, DiffRequest } from "../api/api_pb";

import { loading } from "../widgets";

import client from "../client";

const fetchDiff = async (params): DiffResponse => {
  const req = new DiffRequest(params);
  return await client.diff(req);
};

const Diff: Component = () => {
  const [searchParams] = useSearchParams();

  const [diff] = createResource(() => {
    for (const param of ["a", "b"]) {
      if (searchParams[param] == undefined) {
        throw `Missing search parameter: ${param}`;
      }
    }

    return {
      OutputHash1: searchParams.a,
      OutputHash2: searchParams.b,
    };
  }, fetchDiff);

  return (
    <>
      <Suspense fallback={loading}>
        <Show when={diff()}>
          <div class="card bg-base-100 shadow-xl m-5 w-full h-screen">
            <div class="card-body h-screen">
              <h2 class="card-title">Diffoscope</h2>

              {(() => {
                const iframe = (
                  <iframe
                    class="w-full h-screen"
                    onLoad={() => {
                      const domDoc =
                        iframe.contentDocument || iframe.contentWindow.document;
                      domDoc.write(diff()?.HTMLDiff);
                    }}
                   />
                );

                return <>{iframe}</>;
              })()}
            </div>
          </div>
        </Show>
      </Suspense>
    </>
  );
};

export default Diff;
