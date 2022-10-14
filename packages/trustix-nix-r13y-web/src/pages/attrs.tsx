import {
  Component,
  createResource,
  For,
  Show,
  Suspense,
} from "solid-js";
import { createStore } from "solid-js/store";

import {
  AttrReproducibilityTimeSeriesGroupedbyChannelRequest,
  AttrReproducibilityTimeSeriesResponse,
  AttrReproducibilityTimeSeriesPoint,
} from "../api/api_pb";
import { NameValuePair } from "../lib";

import { loading } from "../widgets";

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
  console.log(timestamps)

  type pointsT = AttrReproducibilityTimeSeriesPoint[];

  const pointsByTimestamp: { [key: number]: pointsT } = { }
  attrKeys
    .forEach((attrKey) => {
      const points = attrs[attrKey].Points

      points.forEach((point) => {
        const ts = Number(point.EvalTimestamp)

        let groupedPoints: pointsT
        if (!(ts in pointsByTimestamp)) {
          groupedPoints = [ ]
          pointsByTimestamp[ts] = groupedPoints
        } else {
          groupedPoints = pointsByTimestamp[ts]
        }

        groupedPoints.push(point)
      })
    })
  console.log(pointsByTimestamp)

    /* .map((attr) => attr.Points.map((p) => p.EvalTimestamp))
     * .map(Number), */

  /* console.log(timestamps); */

    /*
       attrKeys.map((attr) => {
     *       const points = attrs[attr].Points;

     *       const pointsByTimestamp: {
     *         [key: number]: AttrReproducibilityTimeSeriesPoint;
     *       } = {};
     *       points.forEach((p) => {
     *         pointsByTimestamp[Number(p.EvalTimestamp)] = p;
     *       });

     *       return {
     *         label: attr,
     *         data: [6, 5],
     *       };
       }), */
    const chartSettings: SolidChartProps = {
          type: "line",
          data: {
            /* labels: timestamps
             *   .map((ts) => ts * 1000)
             *   .map((ts) => new Date(ts).toISOString()), */
            labels: ["1", "2", "3", "4"],
            datasets: [
              {
                label: "A",
                data: [50, 75, 50, 10],
              },
              {
                label: "B",
                data: [20, 15, 80, 30],
              },
              {
                label: "C",
                data: [80, 30, 80, 19],
                spanGaps: true,
              }
            ],
          },
          options: {
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
