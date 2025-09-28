import React, { useState, useEffect, useCallback } from 'react';

// Wails
import { EventsOnce } from '../../../wailsjs/runtime';
import { server } from "../../../wailsjs/go/models";
import { GetCounters, SelectFile } from "../../../wailsjs/go/backend//App";

// Components
import Grid from '../../components/Grid';
import Header from './Header';

// React component
export default function App() {
    const [path, setPath] = useState<string>("No file selected");
    const [counters, setCounters] = useState<server.CounterImage[]>([]);
    const [filter, setFilter] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(true);
    const [curProgress, setCurProgress] = useState<number>(0);
    const [totalCounters, setTotalCounters] = useState<number>(0);

    const showFolderDialog = useCallback(() => {
        setPath("");
        setLoading(true);
        setCounters([]);
        setCurProgress(0);
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
        if (path === " No file selected ") {
            return (<h2>Select a file to visualize</h2>)
        }

        if (loading) {
            const fileMsg = path === "" ? "Select file..." : `Loading file "${path}"`;
            return (
                <div className='w-full flex flex-col items-center justify-center px-[1ch]'>
                    <p className='text-inherit'>{fileMsg}</p>
                    <progress
                        className="w-full progress"
                        value={curProgress}
                        max="100" />
                </div>
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

        if (curProgress === 0) {
            setTotalCounters(countersLeft);
        }

        const percent = countersLeft / totalCounters;

        setCurProgress(100 - (percent * 100));
    });

    return (
        <div className="flex flex-col items-center">
            <div className='flex flex-col w-[80%] pt-[1ch] gap-[1ch]'>
                <Header
                    filter={filter}
                    setFilter={setFilter}
                    path={path}
                    showFolderDialog={showFolderDialog} />
                <div className='px-[2ch]'>
                    {content()}
                </div>
            </div>
        </div>
    )
}
