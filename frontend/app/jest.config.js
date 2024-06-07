module.exports = {
  resetMocks: true,
  restoreMocks: true,
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!**/*.d.ts',
    '!**/node_modules/**',
    '!src/theme/**',
    '!**/types/**',
  ],
  moduleNameMapper: {
    // Handle CSS imports (without CSS modules)
    '^.+\\.(css|scss)$': '<rootDir>/node_modules/jest-css-modules',

    // Handle module aliases
    '^@/(.*)$': '<rootDir>/src/$1',
    '@redux/(.*)$': '<rootDir>/src/redux/$1',
    '@api/(.*)$': '<rootDir>/src/api/$1',
    '@components/(.*)$': '<rootDir>/src/components/$1',
    '@pages/(.*)$': '<rootDir>/src/pages/$1',
    '@commons/(.*)$': '<rootDir>/src/types/$1',
    '@utils/(.*)$': '<rootDir>/src/utils/$1',
    '@app/(.*)$': '<rootDir>/src/app/$1',
    '@hooks/(.*)$': '<rootDir>/src/hooks/$1',
    '@shared/(.*)$': '<rootDir>/src/shared/$1',
    '@helpers/(.*)$': '<rootDir>/src/helpers/$1',
  },
  transform: {
    '^.+\\.(ts|tsx)?$': 'ts-jest',
  },
  testPathIgnorePatterns: ['<rootDir>/node_modules/'],
  testEnvironment: 'jsdom',
  transformIgnorePatterns: ['^.+\\.module\\.(css|sass|scss)$'],
};
