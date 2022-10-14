import {
  Component,
  createResource,
  For,
  Show,
  Suspense,
  createSignal,
} from "solid-js";
import { createStore } from "solid-js/store";
import { Navigate } from "@solidjs/router";

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

  const datasets = attrKeys.map((attrKey) => {
    const points = timestamps.map((ts) => pointsByTimestamp[ts][attrKey])

    return {
      label: attrKey,
      data: points.map(point => point == undefined ? point : point.PctReproduced),
      "x-r13y-drv": points.map(point => point == undefined ? point : point.DrvPath),
      backgroundColor: palette('tol', points.length).map(hex => `#${hex}`),
      spanGaps: true,
    };
  });

  const [redirStorePath, setRedirStorePath] = createSignal();

  const chartSettings: SolidChartProps = {
    type: "line",
    data: {
      labels: labels,
      datasets: datasets,
    },
    options: {
      onClick: (event, elements) => {
        const drvPaths = new Set<string>();

        for (const elem of elements) {
          const i = elem.index;
          const dsi = elem.datasetIndex;

          const drvPath = datasets[dsi]["x-r13y-drv"][i];
          drvPaths.add(drvPath);
        }

        switch (drvPaths.size) {
          case 0:
            return;
          case 1:
            setRedirStorePath(Array.from(drvPaths)[0]);
            return;
          default:
            alert("multiple derivations at selection point");
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

  return (
    <>
      <h2 className="text-xl font-bold text-center mb-2">{channel}</h2>

      <Show when={redirStorePath()}>
        <Navigate
          href={`/drv?storePath=${encodeURIComponent(redirStorePath())}`}
        />
      </Show>

      <SolidChart
        {...chart}
        canvasOptions={{
          width: 900,
          height: 300,
        }}
      />

      <div className="divider"></div>
    </>
  );
};

const renderChannels = (resp: DerivationReproducibilityResponse): Component => {
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
