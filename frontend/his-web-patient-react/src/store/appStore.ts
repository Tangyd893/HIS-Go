import { create } from 'zustand'

interface AppState {
  collapsed: boolean
  loading: boolean
  toggleCollapsed: () => void
  setLoading: (val: boolean) => void
}

export const useAppStore = create<AppState>((set) => ({
  collapsed: false,
  loading: false,
  toggleCollapsed: () => set((s) => ({ collapsed: !s.collapsed })),
  setLoading: (val: boolean) => set({ loading: val }),
}))
