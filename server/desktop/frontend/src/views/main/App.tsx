import React, { useState, useEffect, useCallback, useMemo, ReactNode, ReactHTMLElement } from 'react';

// Wails
import { EventsOn, EventsOnce } from '../../../wailsjs/runtime';
import { server } from "../../../wailsjs/go/models";
import { GetCounters, SelectFile } from "../../../wailsjs/go/main/App";

// Components
import Grid from '../../components/Grid';
import Header from './Header';

// React component
export default function App() {
    const [path, setPath] = useState<string>("No file selected");
    const [counters, setCounters] = useState<server.CounterImage[]>([]);
    const [filter, setFilter] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(true);
    const [progress, setProgress] = useState<number>(0);
    const [totalCounters, setTotalCounters] = useState<number>(0);

    const showFolderDialog = useCallback(() => {
        setPath("");
        setLoading(true);
        setCounters([]);
        setProgress(0);
        SelectFile().
            then(setPath);
    }, []);

    // Use this to test the counters layout
    // const [counters, setCounters] = useState(() => {
    //     let cs = [];
    //     let counter = {
    //         id: 'example',
    //         pretty_name: 'Wagner_Ground Forces_Infantry_back.png',
    //         filename: counterImg,
    //     };
    //     for (let i = 0; i < 24; i++) {
    //         cs.push(counter);
    //     }
    //     return cs;
    // });

    const content = () => {
        if (path === "No file selected") {
            return (
                <h1 className="flex-initial p-2 mx-2 text-2xl content-center text-center">
                    Select a file to visualize
                </h1>
            )
        }

        if (loading) {
            const fileMsg = path === "" ? "Select file..." : `Loading file "${path}"`;
            return (
                <h1 className="flex flex-col p-2 m-2 text-2xl items-center">
                    <div className='m-2'>{fileMsg}</div>
                    <progress
                        className="progress m-2 w-1/2"
                        value={progress}
                        max="100" />
                </h1>
            )
        }

        return (
            <Grid counters={
                counters.filter((counter) =>
                    counter.pretty_name.toLowerCase()
                        .includes(filter.toLowerCase()))
            } />
        )
    };

    useEffect(() => {
        EventsOnce("counters", async (data?: any) => {
            const counters = await GetCounters();
            setFilter("");
            setLoading(false);
            setCounters(counters);
        })
    }, [path]);

    EventsOnce("processed_left", async (countersLeft?: any) => {
        if (countersLeft === Infinity ||
            countersLeft === 'undefined' ||
            countersLeft === undefined ||
            totalCounters === undefined) {
            return;
        }

        if (progress === 0) {
            setTotalCounters(countersLeft);
        }

        const percent = countersLeft / totalCounters;

        setProgress(100 - (percent * 100));
    });

    return (
        <div className='flex flex-col h-screen'>
            <Header
                filter={filter}
                setFilter={setFilter}
                path={path}
                showFolderDialog={showFolderDialog} />
            {content()}
        </div>
    )
}
