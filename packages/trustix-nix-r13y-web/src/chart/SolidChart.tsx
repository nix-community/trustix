import { Chart, registerables } from "chart.js";
import { ComponentProps, createEffect, createSignal } from "solid-js";
import merge from "lodash.merge";
import deepClone from "lodash.clonedeep";

export type SolidChartProps = {
  canvasOptions?: ComponentProps<"canvas">;
  replace?: boolean; // set this to true if merging properties is not suitable (i.e. when you need to remove items or properties)
} & Chart.ChartConfiguration;

Chart.register(...registerables);

const replaceChartProps = (props: { src: SolidChartProps; dest: Chart }) => {
  const { src, dest } = props;
  for (const key in src) {
    if (key in dest) {
      dest[key] = src[key];
    }
  }
};

export function SolidChart(props: SolidChartProps) {
  const [canvas, setCanvas] = createSignal<HTMLCanvasElement | null>(null);
  const [chart, setChart] = createSignal<Chart | null>(null);
  createEffect(() => {
    const el = canvas();
    if (!el) return;
    const _chart = chart();
    if (!_chart) {
      setChart(new Chart(el, deepClone(props)));
      return;
    }

    if (props.replace) {
      replaceChartProps({ src: props, dest: _chart.config });
    } else {
      merge(_chart.config, deepClone(props));
    }
    _chart.update();
  });
  function createChart(canvas: HTMLCanvasElement) {
    setTimeout(() => setCanvas(canvas));
  }
  return <canvas {...props.canvasOptions} ref={createChart} />;
}
