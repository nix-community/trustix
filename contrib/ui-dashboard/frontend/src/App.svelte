<script>
	export let channels;
	export let evaluations;

    import 'bulma/css/bulma.css'
    import '@fortawesome/fontawesome-free/css/all.css'

    import { onMount } from 'svelte';
	import * as d3 from 'd3';

    import { Input } from 'svelma'
    import { Select } from 'svelma'
    import { Field } from 'svelma'
    import { Button } from 'svelma'


    var data = [30, 86, 168, 281, 303, 365];
	let el;

    onMount(() => {
		d3.select(el)
			.selectAll("div")
			.data(data)
			.enter()
			.append("div")
			.style("width", function(d) {
				return d + "px";
			})
			.text(function(d) {
				return d;
			});
	});

</script>

<main>
    <div class="flex-container">

      <h2>Trustix (TODO: Logo) - &nbsp;</h2>

      <Input type="text" placeholder="Attribute" />

      <Select placeholder="Channel">
        {#each channels as channel}
          <option value={channel}>{channel}</option>
        {/each}
      </Select>

      <Select placeholder="Evaluation">
        {#each evaluations as evaluation}
          <option value={evaluation}>{evaluation}</option>
        {/each}
      </Select>

      <button type="is-primary">Submit</button>

    </div>

    <hr />

    <div class="flex-container">
      <div bind:this={el} class="chart"></div>
    </div>

    <div class="flex-container">

      <div>
        <ul>
          <li><strong>Attribute:</strong> hello.x86_64-linux</li>
          <li><strong>Channel:</strong> {channels[0]}</li>
          <li><strong>Evaluation:</strong> {channels[0]}</li>
        </ul>
      </div>

    </div>


</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

    .flex-container {
        display: flex;
    }

	.chart :global(div) {
		font: 10px sans-serif;
		background-color: steelblue;
		text-align: right;
		padding: 3px;
		margin: 1px;
		color: white;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
