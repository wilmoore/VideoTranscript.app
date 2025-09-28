# Feature: Transcript Display and Downloads

## Overview
Implement transcript viewing and download functionality for the VideoTranscript.app dashboard.

## Problem Statement
1. Job detail page shows metadata but no transcript content
2. Download functionality is a placeholder alert
3. Dashboard uses separate job tracking system (jobs.json) that lacks transcript data
4. Jobs fail due to missing whisper.cpp dependency

## Requirements
- **Long-term persistence**: Use Encore database + Supabase if needed
- **Format priority**: Transcript text → SRT → other formats
- **UI approach**: Simple, user-friendly transcript display
- **Security**: No restrictions for now
- **Fix strategy**: Easiest fixes first

## Technical Analysis

### Current Architecture
- **Dashboard**: `web-dashboard.go` - Frontend with jobs.json tracking
- **Encore Service**: `transcribe/service.go` - Proper database persistence
- **Job Models**: Different structures between systems
- **Missing Dependency**: whisper.cpp not installed

### Root Causes
1. **Architecture Gap**: Dashboard doesn't integrate with Encore's persistent storage
2. **Missing Tools**: whisper.cpp dependency missing
3. **Data Disconnect**: Two separate job tracking systems

## Implementation Plan

### Phase 1: Fix Job Processing (Easiest First)
1. Install whisper.cpp dependency
2. Connect dashboard to Encore database
3. Update job status tracking

### Phase 2: Enhanced Transcript Display
1. Add transcript section to job detail page
2. Implement timestamped segments view
3. Add search and copy functionality

### Phase 3: Download Implementation
1. Create download endpoints for multiple formats
2. Update UI with download buttons
3. File serving with proper MIME types

### Phase 4: Storage Integration
1. Leverage Encore's database for persistence
2. File management for subtitle files
3. API integration between systems

## Files to Modify
- `web-dashboard.go` - Add transcript display and downloads
- `transcribe/service.go` - Potentially add download endpoints
- `transcribe/db.go` - Database queries for dashboard
- Template sections in `transactionDetailHTML`

## Success Criteria
- Job processing works without failures
- Transcript content displays on job detail page
- Download functionality works for all formats (txt, srt, vtt, json)
- Dashboard uses persistent Encore database
- No errors, bugs, or warnings