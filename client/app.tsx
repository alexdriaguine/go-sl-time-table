import { type Component } from "solid-js";

export const App: Component = () => {
	return (
		<section>
			<h1>App for SL time table</h1>
			<button
				onClick={() => {
					fetch("/api/sites?term=sundbyberg")
						.then((r) => r.json())
						.then(console.log);
				}}
			>
				Click me
			</button>
		</section>
	);
};
