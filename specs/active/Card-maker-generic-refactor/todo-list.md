# Todo List: Card-maker-generic-refactor

**Task ID:** Card-maker-generic-refactor  
**Based on:** plan.md Phase 4a, 4b

## Phases

### Phase 4a: Effects Usage in Templates
- [x] Add effectLookup FuncMap helper and wire into render.go
- [x] Update WeaponPart.html to use effect tooltips (title with description)
- [x] Update Item.html to use effect tooltips (title with description)

### Phase 4b: Complete Data Migration
- [x] Complete weapon_parts.yaml (6→25 entries from WeaponPart.csv)
- [x] Complete items.yaml (7→13 entries from Item.csv)

## Progress Log

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | effectLookup FuncMap | done | buildTemplateFuncMap in render.go |
| 2 | WeaponPart.html effects | done | title="{{effectLookup .}}" |
| 3 | Item.html effects | done | title="{{effectLookup .}}" |
| 4 | weapon_parts.yaml | done | 25 entries |
| 5 | items.yaml | done | 13 entries |
