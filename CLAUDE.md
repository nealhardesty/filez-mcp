Enclosed are a set of principles to be appled to all Claude Code Projects


---
# General Software Development Principles

## **Core Development Philosophy**

### Research & Understanding First
- **ALWAYS** read through relevant documentation and specifications before starting work
- Think hard, take time, be thorough - do not rush implementation
- Do your research and make sure any documentation references are loaded, read, and thoroughly understood
- Stay focused and maintain project scope boundaries
- Do not allow for hallucinations.  If a rule, spec, prompt, or documentation is not understood perfectly, ask the user for more information.

### Software Development Principles
- **KISS**: Keep It Simple (when possible) - prefer straightforward solutions over complex ones
- **DRY**: Don't Repeat Yourself - extract common functionality and avoid code duplication
- **Documentation**: Keep comments and README.md up to date ALWAYS - documentation is key
- Be creative within bounds, but don't go off the "deep end"

### Tests
- Implement unit tests **ALWAYS***
- Implement functional tests when it makes sense
- Always create a CLI helper test script to exercise functionality from the command line.  Create this under a scripts/ directory appropriate for the project.


### Makefile
- Makefiles for common build/run/test steps is always helpful

### Quality & User Experience
- UI should look good to humans but prioritize being quick and functional
- Write clean, maintainable code that others can understand
- Test thoroughly with proper error handling
- Prioritize understanding over speed of implementation

### Code Organization
- Maintain clear separation of concerns
- Extract common functionality into reusable components
- Keep functions focused and under reasonable length
- Use meaningful names that clearly describe purpose

### Documentation Standards
- Keep inline comments up to date with code changes
- Important: Maintain README.md files with current setup and usage instructions (create them if it doesn't exist)
- Document complex logic and business rules clearly
- Include examples in documentation when helpful

---
# TAO Documentation Protocol

## **MANDATORY: Thought, Action, Observation Documentation**

It is critical to track all of the actions, thought process and observations made when generating or modifying code.  You will need to create a TAO file for each step of each SPEC file you read in.

### TAO File Creation Requirements
- **MUST** create corresponding `TAO-X.md` file after completing each `SPEC-X.md` task
- TAO files document **Thoughts, Actions, and Observations** chronologically
- Place TAO files in same directory as the `SPEC-X.md` file
- Each TAO file serves as a complete log of the implementation process.  I do mean COMPLETE!
- After all of the steps have been documented and written, create a new section that is a concise summary of all actions taken, similar to the overview displayed in the agent summary.

### Required TAO Format (EXACT FORMAT REQUIRED):
```markdown
### Step [Step Number]
#### Thought
[Your internal reasoning here. Explain why you're taking this step, what you're trying to achieve, and what tools or files you plan to use. Be concise but thorough]

#### Action
[The specific action you're taking. Be concise and use clear format: Tool: [Tool Name], File: [File Name], or Command: [Terminal Command]. List files modified/created/deleted]

#### Observation
[The result of your action. Tool output, file content, or command result. This data will inform your next thought.]
```

### TAO Documentation Rules
- Document EVERY step in the TAO format
- Be thorough and transparent in documentation
- **Do NOT proceed to next step until current Observation is fully processed**
- Each step must have complete Thought → Action → Observation cycle
- IMPORTANT, when handling ad-hoc prompts from the user, you MUST also update the latest TAO-X.md file with the same updates you would for processing a SPEC file.
- Final output should include completed TAO file and modified code files

### TAO Process Guidelines
- Use sequential step numbering (Step 1, Step 2, etc.)
- Each Thought should explain the reasoning behind the action
- Each Action should clearly state what is being done and to which files
- Each Observation should capture the complete result/output
- Use the observation data to inform the next thought process

### TAO Quality Standards
- Be concise but thorough in explanations
- Include specific file names and tool names in actions
- Capture complete output in observations (don't truncate important details)
- Maintain chronological order throughout the process
- End each TAO file with a summary of what was accomplished
