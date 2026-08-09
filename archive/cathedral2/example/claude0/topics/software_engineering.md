# Software Engineering

## Topic Overview

Professional domain knowledge related to software development, with focus on [[people/alex_chen]]'s expertise in frontend development, React ecosystem, and engineering practices.

## Tags
`professional` `technical` `react` `frontend` `programming`

## Technical Stack

### Primary Technologies

#### Frontend Focus
- **React**: Primary framework (5+ years experience)
- **TypeScript**: Strongly preferred over JavaScript
- **State Management**: Redux, Zustand, Context API
- **Styling**: CSS-in-JS, Tailwind CSS
- **Testing**: Jest, React Testing Library, Cypress

#### Full-Stack Capabilities
- **Backend**: Python (FastAPI, Django)
- **Databases**: PostgreSQL, Redis
- **DevOps**: Docker, K8s, GitHub Actions
- **Cloud**: AWS (primary), some GCP

### Development Environment

Alex's setup:
- **Editor**: VS Code with Vim keybindings
- **Terminal**: iTerm2 with Oh My Zsh
- **Version Control**: Git (obviously)
- **Machine**: MacBook Pro M1 Max

> "I spend way too much time customizing my dev environment. It's like gear acquisition syndrome but for dotfiles." - [[episodes/2024-02-18_debugging_session]]

## Engineering Philosophy

### Code Quality Principles

1. **Functional Approach**: Prefers pure functions, immutability
2. **Testing**: "If it's not tested, it's broken"
3. **Documentation**: Code should be self-documenting, but document the why
4. **Performance**: Measure first, optimize second

### Problem-Solving Approach

Demonstrated in [[episodes/2024-02-18_debugging_session]]:
- Systematic debugging (eliminate variables)
- Root cause analysis
- Production vs. development awareness
- Tool utilization (React DevTools, profilers)

## Professional Context

### Current Role
- **Title**: Senior Software Engineer
- **Company**: TechCorp (3 years)
- **Team**: Lead for 4 developers
- **Focus**: E-commerce platform, checkout flow

### Career Trajectory
1. Started at startup (2 years) - Full-stack generalist
2. Second startup (1.5 years) - Frontend specialization
3. Current role - Technical leadership, architecture

### Work Style
- Hybrid (2 days office/week)
- Morning person (best coding before 10 AM)
- Advocates for work-life balance
- Regular "mental health deployments" [[topics/hiking_pacific_northwest]]

## Technical Challenges Encountered

### Recent Issues Solved

1. **React Performance** (First meeting):
   - Component re-rendering 50 times/second
   - Solution: memo and useCallback optimization
   - Result: 95% reduction in renders

2. **Hooks Production Bug** [[episodes/2024-02-18_debugging_session]]:
   - useEffect infinite loop in production only
   - Root cause: Build optimization changing function references
   - Solution: useRef for stable references

### Lessons Learned

> "I've been thinking about hooks wrong. I was treating them like lifecycle methods instead of synchronization primitives." - Alex

Key insights:
- Production builds can behave differently
- Less complexity often better than more
- Understanding fundamentals crucial

## Cross-Domain Thinking

### Programming ↔ Music Production

Alex draws parallels [[episodes/2024-01-22_music_production]]:
- Mixing as "debugging audio"
- Arrangement as "system architecture"
- DRY principle doesn't apply to art
- Refactoring concepts in creative work

### Programming ↔ Hiking

Mental models:
- Trail planning like project planning
- Gear optimization like code optimization
- "Clearing cache" through nature
- Systematic approach to both

## Learning & Growth

### Current Learning
- **Rust**: Exploring systems programming
- **WebAssembly**: For performance-critical features
- **AI/ML**: Curious but not actively pursuing

### Teaching & Mentoring
- Leads team learning sessions
- Documents solutions for team wiki
- Patient with junior developers
- Values knowledge sharing

## Tool Preferences

### Development Tools
- **Bundler**: Vite (moving away from Webpack)
- **Package Manager**: pnpm
- **Code Quality**: ESLint, Prettier, Husky
- **Monitoring**: Sentry, DataDog

### Productivity Tools
- Linear for issue tracking
- Notion for documentation
- Excalidraw for architecture diagrams
- Raycast for macOS productivity

## Memorable Quotes

On debugging:
> "Even the rubber duck is judging me" [[episodes/2024-02-18_debugging_session]]

On work-life balance:
> "I spend my days debugging code and my evenings debugging why my synthesizer won't sync with Ableton" [[episodes/2024-01-22_music_production]]

On continuous learning:
> "The moment you think you know everything about React, they release a new hook that breaks your mental model"

## Related Memories

- [[people/alex_chen]] - Personal profile
- [[concepts/react_hooks]] - Deep technical knowledge
- [[topics/python_development]] - Backend skills
- [[episodes/2024-02-18_debugging_session]] - Problem-solving example

---
*Knowledge domain established: 2024-01-15*
*Last updated: 2024-02-18*
*Strength: 0.93*