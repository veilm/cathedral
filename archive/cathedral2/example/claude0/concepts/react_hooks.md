# React Hooks

## Concept Overview

React's system for using state and other React features in functional components. Deep technical knowledge gained through debugging sessions with [[people/alex_chen]].

## Tags
`react` `javascript` `frontend` `programming-concept` `state-management`

## Core Understanding

### Mental Model Evolution

Initial misconception (Alex's):
> "I was treating them like lifecycle methods instead of synchronization primitives" [[episodes/2024-02-18_debugging_session]]

Correct understanding:
- Hooks synchronize React with external systems
- Not lifecycle replacements but effect descriptions
- Order matters (Rules of Hooks)
- Closures and dependencies crucial

## Common Hooks Deep Dive

### useState
- Local component state
- Returns [value, setter] tuple
- Setter can take function for previous state
- Batching behavior important

### useEffect

Critical learnings from debugging:
```javascript
// Problematic pattern
useEffect(() => {
  setSocket(new WebSocket(url));
}, [url, setSocket]); // setSocket causes issues

// Better approach
useEffect(() => {
  const ws = new WebSocket(url);
  socketRef.current = ws;
  return () => ws.close();
}, [url]);
```

Key insights:
- Cleanup functions crucial
- Dependency array precision matters
- Production vs. development differences
- Avoid state setters in dependencies

### useRef

Usage patterns discovered:
- Stable reference across renders
- DOM element access
- Mutable values without re-renders
- WebSocket/interval storage

### Custom Hooks

Example from Alex's work:
```javascript
const useOrderTracking = (orderId) => {
  // Complex state and effect management
  // Reusable across components
  // Encapsulates business logic
};
```

## Advanced Patterns

### Performance Optimization

Techniques used successfully:
1. **useMemo**: Expensive calculations
2. **useCallback**: Stable function references
3. **React.memo**: Component memoization
4. **useReducer**: Complex state logic

### Common Pitfalls

From real debugging sessions:

1. **Infinite Loops**:
   - Missing dependencies
   - Creating new objects in render
   - State updates in effects

2. **Stale Closures**:
   - Not including values in dependencies
   - Using state instead of refs for values

3. **Production Build Issues**:
   - Minimization changing behavior
   - Double-rendering in StrictMode
   - Different optimization paths

## Debugging Strategies

### Tools and Techniques

Successful approaches:
- React DevTools Profiler
- Console logging render counts
- Breakpoints in effects
- Dependency array analysis

### Mental Debugging Process

Alex's evolved approach:
1. Identify unexpected behavior
2. Check dependency arrays
3. Look for new object creation
4. Consider closure issues
5. Test in production build

## Real-World Examples

### The WebSocket Problem

From [[episodes/2024-02-18_debugging_session]]:

**Problem**: Infinite re-renders in production
**Root Cause**: Function reference instability
**Solution**: useRef for stable reference
**Lesson**: Production optimizations matter

### Performance Optimization

First meeting issue:
- 50 re-renders per second
- Solution: memo + useCallback
- Result: 95% reduction
- Key: Understanding render triggers

## Best Practices Developed

### Alex's Rules

1. "Minimal dependency arrays"
2. "Refs for non-render values"
3. "Effects for synchronization only"
4. "Test in production build"

### Code Organization

Patterns adopted:
- Custom hooks for complex logic
- Separate concerns clearly
- Name hooks descriptively
- Document unusual patterns

## Learning Resources

### Mentioned References
- React docs (new beta version)
- Kent C. Dodds articles
- Dan Abramov's blog
- React DevTools tutorial

### Community Wisdom
- "Hooks are hard until they click"
- "Think synchronization, not lifecycle"
- "When in doubt, check dependencies"

## Evolution & Future

### React's Direction
- Server Components impact
- Concurrent features
- Suspense integration
- Compiler optimizations

### Alex's Perspective
> "The moment you think you know everything about React, they release a new hook that breaks your mental model"

## Common Patterns Library

Built from experience:
- Data fetching hooks
- WebSocket management
- Form handling
- Animation hooks
- Local storage sync

## Related Memories

- [[topics/software_engineering]] - Professional context
- [[people/alex_chen]] - User profile
- [[episodes/2024-02-18_debugging_session]] - Debugging session
- [[episodes/2024-01-15_first_meeting]] - Initial React discussion

---
*Concept documented: 2024-01-15*
*Last updated: 2024-02-18*
*Strength: 0.92*