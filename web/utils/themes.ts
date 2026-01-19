export type ThemeKey =
  | 'newyear'
  | 'birthday'
  | 'promotion'
  | 'anniversary'
  | 'memorial'
  | 'proposal'
  | 'baby'
  | 'admission'
  | 'religious'
  | 'custom'

export const THEMES: Record<ThemeKey, { label: string; bg: string; primary: string; secondary: string }> = {
  newyear: {
    label: 'Новый год',
    bg: 'bg-gradient-to-br from-sky-900 via-slate-950 to-indigo-950',
    primary: '#38bdf8',
    secondary: '#a78bfa',
  },
  birthday: {
    label: 'День рождения',
    bg: 'bg-gradient-to-br from-fuchsia-700 via-slate-950 to-amber-600',
    primary: '#d946ef',
    secondary: '#f59e0b',
  },
  promotion: {
    label: 'Повышение',
    bg: 'bg-gradient-to-br from-indigo-700 via-slate-950 to-sky-600',
    primary: '#6366f1',
    secondary: '#0ea5e9',
  },
  anniversary: {
    label: 'Юбилей',
    bg: 'bg-gradient-to-br from-rose-600 via-slate-950 to-orange-600',
    primary: '#fb7185',
    secondary: '#fb923c',
  },
  memorial: {
    label: 'Памятная дата',
    bg: 'bg-gradient-to-br from-slate-800 via-slate-950 to-slate-700',
    primary: '#cbd5e1',
    secondary: '#94a3b8',
  },
  proposal: {
    label: 'Предложение',
    bg: 'bg-gradient-to-br from-pink-600 via-slate-950 to-rose-600',
    primary: '#f472b6',
    secondary: '#fb7185',
  },
  baby: {
    label: 'Рождение ребёнка',
    bg: 'bg-gradient-to-br from-sky-500 via-slate-950 to-teal-500',
    primary: '#38bdf8',
    secondary: '#2dd4bf',
  },
  admission: {
    label: 'С поступлением',
    bg: 'bg-gradient-to-br from-lime-500 via-slate-950 to-emerald-600',
    primary: '#84cc16',
    secondary: '#10b981',
  },
  religious: {
    label: 'Религиозные',
    bg: 'bg-gradient-to-br from-amber-500 via-slate-950 to-yellow-500',
    primary: '#f59e0b',
    secondary: '#eab308',
  },
  custom: {
    label: 'Кастомный',
    bg: 'bg-slate-950',
    primary: '#a78bfa',
    secondary: '#22d3ee',
  },
}
