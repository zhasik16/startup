package handlers

import (
    "testing"
)

func TestAutoFixEngine(t *testing.T) {
    engine := NewAutoFixEngine()
    
    testRisks := []Risk{
        {
            Title:       "Hardcoded API Key",
            CodeSnippet: `API_KEY = "sk-1234567890abcdef"`,
        },
        {
            Title:       "Hardcoded Password", 
            CodeSnippet: `password = "super_secret_123"`,
        },
    }
    
    testCodebase := map[string]string{
        "app.py": `API_KEY = "sk-1234567890abcdef"
password = "super_secret_123"`,
    }
    
    fixes := engine.GenerateFixes(testRisks, testCodebase)
    
    if len(fixes) == 0 {
        t.Error("Expected auto-fixes to be generated")
    }
    
    for _, fix := range fixes {
        if fix.Original == fix.Fixed {
            t.Errorf("Fix didn't change code: %s", fix.RiskTitle)
        }
        if fix.Explanation == "" {
            t.Errorf("Missing explanation for: %s", fix.RiskTitle)
        }
    }
}