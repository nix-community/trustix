// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import {
  Component,
  createResource,
  For,
  Show,
  Suspense,
  createSignal,
} from "solid-js";
import { createStore } from "solid-js/store";
import { Navigate, A } from "@solidjs/router";

import {
  AttrReproducibilityTimeSeriesGroupedbyChannelRequest,
  AttrReproducibilityTimeSeriesResponse,
  AttrReproducibilityTimeSeriesPoint,
} from "../api/api_pb";
import { NameValuePair } from "../lib";

import { loading } from "../widgets";

import palette from "google-palette";

import { SolidChart, SolidChartProps } from "../chart/SolidChart";

import client from "../client";

const fetchAttrsByChannel = async (): DerivationReproducibilityResponse => {
  const req = new AttrReproducibilityTimeSeriesGroupedbyChannelRequest({});
  return await client.attrReproducibilityTimeSeriesGroupedbyChannel(req);
};

const drvLink = (storePath): string =>
  `/drv?storePath=${encodeURIComponent(encodeURIComponent(storePath))}`;

/* eslint-disable sonarjs/cognitive-complexity */
const renderChannel = (
  channel: string,
  attrs: { [key: string]: AttrReproducibilityTimeSeriesResponse },
): Component => {
  const attrKeys = Object.keys(attrs);

  const timestamps = [
    ...new Set(
      attrKeys
        .map((attrKey) => attrs[attrKey])
        .map((attr) => attr.Points.map((p) => p.EvalTimestamp))
        .map(Number),
    ),
  ];

  type pointT = AttrReproducibilityTimeSeriesPoint;

  const pointsByTimestamp: {
    [key: number]: { [key: string]: AttrReproducibilityTimeSeriesPoint };
  } = {};
  attrKeys.forEach((attrKey) => {
    const points = attrs[attrKey].Points;

    points.forEach((point) => {
      const ts = Number(point.EvalTimestamp);

      let byAttr: { [key: string]: pointT };
      if (ts in pointsByTimestamp) {
        byAttr = pointsByTimestamp[ts];
      } else {
        byAttr = {};
        pointsByTimestamp[ts] = byAttr;
      }

      if (!(attrKey in byAttr)) {
        byAttr[attrKey] = point;
      }
    });
  });

  /* const labels =  */
  const labels = timestamps
    .map((ts) => ts * 1000)
    .map((ts) => new Date(ts).toISOString());

  // If there is only one label make it both the first and the last time in chart
  // so it looks less empty
  if (labels.length == 1) {
    labels.push(labels[0]);
  }

  const chartColours = palette("tol", attrKeys.length).map((hex) => `#${hex}`);

  const datasets = attrKeys.map((attrKey, i) => {
    const points = timestamps.map((ts) => pointsByTimestamp[ts][attrKey]);
    return {
      label: attrKey,
      data: points.map((point) =>
        point == undefined ? point : point.PctReproduced,
      ),
      "x-r13y-drv": points.map((point) =>
        point == undefined ? point : point.DrvPath,
      ),
      backgroundColor: chartColours[i],
      borderColor: chartColours[i],
      spanGaps: true,
    };
  });

  const [redirStorePath, setRedirStorePath] = createSignal();

  const [multiChartSelection, setMultiChartSelection] = createSignal([]);
  const multiChartSelectionID = `multi-point-chart-selection-${channel}`;
  const multiChartSelectionInput = (
    <input
      type="checkbox"
      id={multiChartSelectionID}
      class="modal-toggle"
    />
  );

  const chartSettings: SolidChartProps = {
    type: "line",
    data: {
      labels: labels,
      datasets: datasets,
    },
    options: {
      onClick: (event, elements) => {
        const drvPaths = new Set<string>();
        const attrs: { [key: string]: string[] } = {};

        for (const elem of elements) {
          const dataset = datasets[elem.datasetIndex];
          const attr = dataset["label"];

          let attrDrvs: string[] = [];
          if (attr in attrs) {
            attrDrvs = attrs[attr];
          } else {
            attrDrvs = [];
            attrs[attr] = attrDrvs;
          }

          const drvPath = dataset["x-r13y-drv"][elem.index];

          drvPaths.add(drvPath);
          attrDrvs.push(drvPath);
        }

        switch (drvPaths.size) {
          case 0:
            return;
          case 1:
            setRedirStorePath(Array.from(drvPaths)[0]);
            return;
          default:
            setMultiChartSelection(NameValuePair.fromMap(attrs));
            multiChartSelectionInput.checked = true;
            break;
        }
      },
      responsive: true,
      plugins: {
        legend: {
          position: "top",
        },
      },
      scales: {
        x: {
          display: true,
          title: {
            display: true,
          },
        },
        y: {
          display: true,
          title: {
            display: true,
            text: "%",
          },
          suggestedMin: 0,
          suggestedMax: 100,
        },
      },
    },
  };

  const [chart] = createStore(chartSettings);

  const attrsList = NameValuePair.fromMap(attrs);

  return (
    <>
      <div class="card bg-base-100 shadow-xl m-5">
        <div class="card-body">
          <h2 class="card-title">{channel}</h2>

          {/* when multiple points on the chart occupy the same space in the chart ask the user to clarify selection */}
          {multiChartSelectionInput}
          <div class="modal">
            <div class="modal-box w-11/12 max-w-5xl">
              <h3 class="text-lg font-bold">
                Multiple derivations found at chart point.
              </h3>

              <table class="table w-full">
                <thead>
                  <tr>
                    <th>Attribute</th>
                    <th>Derivations</th>
                  </tr>
                </thead>
                <tbody>
                  <For each={multiChartSelection()}>
                    {({ name, value }) => {
                      const attr: string = name;
                      const drvs: string[] = value;

                      return (
                        <tr>
                          <td>{attr}</td>
                          <td>
                            <For each={drvs}>
                              {(drvPath) => (
                                <>
                                  <A href={drvLink(drvPath)}>{drvPath}</A>
                                </>
                              )}
                            </For>
                          </td>
                        </tr>
                      );
                    }}
                  </For>
                </tbody>
              </table>

              <div class="modal-action">
                <label for={multiChartSelectionID} class="btn">
                  Close
                </label>
              </div>
            </div>
          </div>

          <Show when={redirStorePath()}>
            <Navigate href={drvLink(redirStorePath())} />
          </Show>

          <SolidChart
            {...chart}
            canvasOptions={{
              width: 900,
              height: 300,
            }}
          />

          <div class="overflow-x-auto">
            <table class="table w-full">
              <thead>
                <tr>
                  <th>Attribute</th>
                  <th>Derivations</th>
                </tr>
              </thead>
              <tbody>
                <For each={attrsList}>
                  {({ name, value }) => {
                    const attrName = name;
                    const points = value.Points;

                    let derivationsText = "No derivations…";
                    if (points.length > 0) {
                      derivationsText = points[0].DrvPath + "…";
                    }

                    return (
                      <>
                        <tr>
                          <td>{attrName}</td>
                          <td>
                            <div
                              tabIndex={0}
                              class="collapse collapse-arrow pl-0"
                            >
                              <input type="checkbox" />

                              <div class="collapse-title pl-0">
                                {derivationsText}
                              </div>

                              <div class="collapse-content">
                                <table class="table w-full">
                                  <thead>
                                    <tr>
                                      <th>Timestamp</th>
                                      <th>Derivation</th>
                                    </tr>
                                  </thead>
                                  <tbody>
                                    <For each={points}>
                                      {(p) => {
                                        const evalDate = new Date(
                                          Number(p.EvalTimestamp) * 1000,
                                        );
                                        const evalDateISO =
                                          evalDate.toISOString();

                                        return (
                                          <tr>
                                            <td>{evalDateISO}</td>
                                            <td>
                                              <A href={drvLink(p.DrvPath)}>
                                                {p.DrvPath}
                                              </A>
                                            </td>
                                          </tr>
                                        );
                                      }}
                                    </For>
                                  </tbody>
                                </table>
                              </div>
                            </div>
                          </td>
                        </tr>
                      </>
                    );
                  }}
                </For>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </>
  );
};

const renderChannels = (resp: DerivationReproducibilityResponse): Component => {
  /* eslint-disable solid/reactivity */
  const channels = NameValuePair.fromMap(resp.Channels);

  return (
    <>
      <For each={channels}>
        {({ name, value }) => renderChannel(name, value.Attrs)}
      </For>
    </>
  );
};

const Attrs: Component = () => {
  const [attrsByChannel] = createResource(fetchAttrsByChannel);

  return (
    <>
      <Suspense fallback={loading}>
        <Show when={attrsByChannel()}>{renderChannels(attrsByChannel())}</Show>
      </Suspense>
    </>
  );
};

export default Attrs;
