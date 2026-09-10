import React from 'react';
import { render } from 'ink-testing-library';
import { Box, Text } from 'ink';

// mock stdout
const width = 80;

function TestApp() {
  const jdText = "Project Role Description : Develop custom software solutions\nto design, code, and enhance components\n";
  const lines = jdText.split('\\n');
  return (
    <Box flexDirection="column" width={width}>
      <Text color="blueBright">
        ╭─ <Text bold>Job Description</Text> {'─'.repeat(Math.max(0, width - 19))}
      </Text>
      <Box flexDirection="column" paddingX={2} paddingY={0}>
        {lines.map((line, i) => (
          <Text key={i}>{line}</Text>
        ))}
      </Box>
      <Text color="blueBright">
        ╰{'─'.repeat(Math.max(0, width - 23))}
        <Text color="white" dimColor> [Ctrl+S to analyse] </Text>
      </Text>
    </Box>
  );
}

const { lastFrame } = render(<TestApp />);
console.log(lastFrame());
