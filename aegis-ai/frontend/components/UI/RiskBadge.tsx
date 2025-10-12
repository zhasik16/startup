interface RiskBadgeProps {
  level: 'critical' | 'high' | 'medium' | 'low';
  count?: number;
}

export default function RiskBadge({ level, count }: RiskBadgeProps) {
  const config = {
    critical: {
      color: 'bg-red-500/20 border-red-500/30 text-red-300',
      icon: '🚨'
    },
    high: {
      color: 'bg-orange-500/20 border-orange-500/30 text-orange-300',
      icon: '⚠️'
    },
    medium: {
      color: 'bg-yellow-500/20 border-yellow-500/30 text-yellow-300',
      icon: '📝'
    },
    low: {
      color: 'bg-blue-500/20 border-blue-500/30 text-blue-300',
      icon: 'ℹ️'
    }
  };

  const { color, icon } = config[level];

  return (
    <span className={`inline-flex items-center space-x-1 px-3 py-1 rounded-full border text-sm font-medium ${color}`}>
      <span>{icon}</span>
      <span className="capitalize">{level}</span>
      {count !== undefined && (
        <span className="bg-white/20 rounded-full px-2 py-0.5 text-xs">
          {count}
        </span>
      )}
    </span>
  );
}