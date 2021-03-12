import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import typescript from "@rollup/plugin-typescript";
import json from '@rollup/plugin-json';
import nodePolyfills from 'rollup-plugin-polyfill-node';
import replace from '@rollup/plugin-replace';
import css from 'rollup-plugin-css-only'
import { terser } from "rollup-plugin-terser";

export default {
  input: ["src/main.ts"],
  output: {
    sourcemap: true,
    dir: "dist",
    format: "iife",
  },
  plugins: [
    nodePolyfills(),
    nodeResolve({
      preferBuiltins: true,
      browser: true,
    }),
    commonjs(),
    json(),
    typescript(),
    replace({
      preventAssignment: true,
      'process.env.NODE_ENV': JSON.stringify('production'),
    }),
    css(),
    terser(),
  ]
};
