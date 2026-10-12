<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Study Print</title>
    @php
        $lineStepPx = round($fontSize * $lineHeight, 2);
        $lineStep = $lineStepPx.'px';
        $ruleLineSvg = 'data:image/svg+xml,'.rawurlencode(
            sprintf(
                '<svg xmlns="http://www.w3.org/2000/svg" width="4" height="%1$s" viewBox="0 0 4 %1$s" preserveAspectRatio="none" shape-rendering="crispEdges">'
                .'<rect x="0" y="%2$s" width="4" height="1" fill="#444"/>'
                .'</svg>',
                $lineStepPx,
                max(0, $lineStepPx - 1),
            )
        );
    @endphp
    <style>
        @page {
            size: A4 {{ $landscape ? 'landscape' : 'portrait' }};
            margin: 12mm;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: {{ $fontFamily }};
            font-size: {{ $fontSize }}px;
            line-height: {{ $lineHeight }};
            text-align: {{ $justifyText ? 'justify' : 'left' }};
            color: #111;
            -webkit-print-color-adjust: exact;
            print-color-adjust: exact;
        }

        .columns {
            display: flex;
            width: 100%;
            min-height: calc(100vh - 24mm);
            align-items: stretch;
        }

        .column {
            display: flex;
            flex-direction: column;
            flex: 1 1 0;
            min-width: 0;
            min-height: calc(100vh - 24mm);
            padding: 0 4mm;
            border-right: 1px solid #ccc;
        }

        .column:last-child {
            border-right: none;
        }

        .column-title {
            font-weight: 600;
            font-size: 0.95em;
            margin-bottom: 0.75em;
            padding-bottom: 0.35em;
            border-bottom: 1px solid #ddd;
            flex-grow: 0;
        }

        .column-subtitle {
            font-size: 0.75em;
            color: #666;
            margin-bottom: 0.75em;
        }

        .column-bible .reader p,
        .column-scribe-content .reader p {
            margin-bottom: 0.75em;
        }

        .reader sup {
            font-size: 0.7em;
            color: #555;
            margin-right: 0.15em;
        }

        .notes-content {
            white-space: pre-wrap;
        }

        .lined-area {
            position: relative;
            flex: 1 1 auto;
            min-height: 0;
        }

        .lined-block {
            position: absolute;
            inset: 0;
            -webkit-print-color-adjust: exact;
            print-color-adjust: exact;
            image-rendering: crisp-edges;
            background-image: url('{!! $ruleLineSvg !!}');
            background-size: 100% {{ $lineStep }};
            background-repeat: repeat-y;
            background-position: top left;
        }

        .column-scribe .lined-area,
        .column-lined-notes .lined-area {
            margin: 0;
        }

        .column-scribe p,
        .column-lined-notes p,
        .column-scribe .reader p,
        .column-lined-notes .reader p {
            margin: 0;
        }

        .message {
            color: #666;
            font-size: 0.9em;
            font-style: italic;
        }

        .verse-list-entry {
            margin-bottom: 0.75em;
        }

        .verse-list-entry strong {
            display: block;
            margin-bottom: 0.25em;
        }
    </style>
</head>
<body>
    <div class="columns">
        @foreach ($columns as $column)
            <div class="column column-{{ $column['kind'] }}">
                <div class="column-title">{{ $column['label'] }}</div>

                @if (($column['kind'] ?? '') === 'bible')
                    <div class="reader">
                        @include('print.partials.chapter-paragraphs', ['verses' => $column['verses']])
                    </div>
                @elseif (($column['kind'] ?? '') === 'notes')
                    <div class="notes-content">{{ $column['content'] }}</div>
                @elseif (($column['kind'] ?? '') === 'lined-notes')
                    <div class="column-subtitle">Notes</div>
                    <div class="lined-area">
                        <div class="lined-block" aria-hidden="true"></div>
                    </div>
                @elseif (($column['kind'] ?? '') === 'scribe')
                    <div class="column-subtitle">Scribe</div>
                    <div class="lined-area">
                        <div class="lined-block" aria-hidden="true"></div>
                    </div>
                @elseif (($column['kind'] ?? '') === 'scribe-content')
                    <div class="reader">
                        @include('print.partials.chapter-paragraphs', [
                            'verses' => $column['verses'],
                            'escapeText' => true,
                        ])
                    </div>
                @elseif (($column['kind'] ?? '') === 'message')
                    <p class="message">{{ $column['message'] }}</p>
                @elseif (($column['kind'] ?? '') === 'verse-list')
                    @foreach ($column['entries'] ?? [] as $entry)
                        <div class="verse-list-entry">
                            <strong>{{ $entry['label'] }}</strong>
                            <p>{{ $entry['text'] }}</p>
                        </div>
                    @endforeach
                @endif
            </div>
        @endforeach
    </div>
</body>
</html>
