import React from 'react';
import { useEffect, useState } from 'react';
import { fetchEventSource, EventSourceMessage } from '@microsoft/fetch-event-source';

import Counter, { CounterData } from './Counter';

const serverBaseURL = "http://localhost:8090/api";

export default function Grid() {
	const [counters, setCounters] = useState([]);

	useEffect(() => {
		const fetchData = async () => {
			await fetchEventSource(`${serverBaseURL}/sse`, {
				method: "GET",
				headers: {
					Accept: "text/event-stream",
				},

				onopen(res: Response): Promise<void> {
					if (res.ok && res.status === 200) {
						console.log("Connection made");
					} else if (
						res.status >= 400 &&
						res.status < 500 &&
						res.status !== 429
					) {
						console.log("Client side error ", res);
					}

					return Promise.resolve();
				},

				onmessage(event: EventSourceMessage): void {
					fetch("/api/counters")
						.then((response) => {
							if (!response.ok) {
								throw new Error(`HTTP error! status: ${response.status}`);
							}
							return response.json();
						}).then((data) => {
							setCounters(data.counters);
						});
				},

				onclose(): void {
					console.log("Connection closed by the server");
				},

				onerror(err): void {
					console.log("There was an error from server", err);
				},
			});
		};
		fetchData();
	}, []);

	return (
		<div className="grid grid-cols-12 gap-2">
			{
				counters.map((counter: CounterData) =>
					<Counter key={counter.id} counter={counter} />
				)
			}
		</div>
	);
}

