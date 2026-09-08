import React from 'react'
import { TerminalFrame } from './components/TerminalFrame/TerminalFrame'
import { Dashboard } from './pages/Dashboard/Dashboard'
import { Send } from './pages/Send/Send'
import { Receivers } from './pages/Receivers/Receivers'
import { Groups } from './pages/Groups/Groups'
import { Templates } from './pages/Templates/Templates'
import { Logs } from './pages/Logs/Logs'
import { Settings } from './pages/Settings/Settings'
import { useUiStore } from './stores/useUiStore'

export const App: React.FC = () => {
  const { activeTab } = useUiStore()

  const renderActivePage = () => {
    switch (activeTab) {
      case 'dashboard':
        return <Dashboard />
      case 'send':
        return <Send />
      case 'receivers':
        return <Receivers />
      case 'groups':
        return <Groups />
      case 'templates':
        return <Templates />
      case 'logs':
        return <Logs />
      case 'settings':
        return <Settings />
      default:
        return <Dashboard />
    }
  }

  return (
    <TerminalFrame>
      {renderActivePage()}
    </TerminalFrame>
  )
}

export default App
