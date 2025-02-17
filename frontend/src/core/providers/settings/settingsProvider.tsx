import { useEffect, useState } from "react";
import { ExperimentalFeatures, Settings, SettingsContext } from "./settingsContext";
import { useLocalStorage } from "@mantine/hooks";

export function SettingsProvider({ loading, children }: { loading: React.ReactNode, children: React.ReactNode }) {
    const [settings, setSettings] = useState<Settings>({} as Settings);
    const [ready, setReady] = useState(false);

    const [experimental, setExperimental] = useLocalStorage<ExperimentalFeatures>({
        key: 'experimental',
        defaultValue: {
            dashboard: false
        }
    })

    useEffect(() => {
        setSettings(prev => ({ ...prev, experimental }))
    }, [experimental])

    useEffect(() => {
        fetch('/settings.json')
            .then((response) => response.json())
            .then((data) => setSettings(prev => ({ ...prev, ...data })))
            .then(() => setReady(true))
            .catch(console.error)
    }, [])

    return (
        <SettingsContext.Provider value={{ settings, setExperimental }}>
            {ready ? children : loading}
        </SettingsContext.Provider>
    )
}