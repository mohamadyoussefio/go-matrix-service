package main

import (
	"fmt"
	"go-matrix-service/internal/protocol"
	"os"
	"os/exec"
	"runtime"
	"text/template"
)

const reportTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Matrix Multiplication | Report</title>
    <link href="https://fonts.googleapis.com/css2?family=Courier+Prime:wght@400;700&display=swap" rel="stylesheet">
    <style>
        * { box-sizing: border-box; }

        body {
            background-color: #f0f0f0; /* Light grey background for contrast */
            color: #000;
            font-family: 'Courier Prime', monospace;
            height: 100vh;
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        /* The "Receipt" Card */
        .card {
            background: #fff;
            width: 600px;
            max-width: 90%;
            border: 3px solid black;
            box-shadow: 15px 15px 0px rgba(0,0,0,0.15);
            padding: 40px;
            display: flex;
            flex-direction: column;
            gap: 30px;
        }

        header {
            border-bottom: 3px solid black;
            padding-bottom: 20px;
            text-align: center;
        }

        h1 { margin: 0; font-size: 24px; text-transform: uppercase; letter-spacing: 2px; }
        .subtitle { color: #666; font-size: 14px; margin-top: 5px; }

        /* The Hero Section (Speedup) */
        .hero { text-align: center; padding: 10px 0; }
        .hero-label { font-size: 14px; color: #666; text-transform: uppercase; margin-bottom: 5px; }
        .hero-val { font-size: 64px; font-weight: bold; line-height: 1; }
        .hero-note { font-size: 12px; color: #666; margin-top: 10px; }

        /* CSS Visualization Bars */
        .viz-container {
            display: flex;
            flex-direction: column;
            gap: 15px;
            padding: 20px;
            background: #fafafa;
            border: 1px dashed black;
        }

        .bar-group { display: flex; align-items: center; gap: 15px; }
        .bar-label { width: 100px; font-size: 12px; font-weight: bold; text-transform: uppercase; text-align: right; }

        .track { flex-grow: 1; height: 24px; background: #eee; position: relative; }
        .fill { height: 100%; background: black; transition: width 1s ease-out; }

        .time-text { margin-left: 10px; font-size: 14px; font-weight: bold; width: 80px; }

        /* Data Table */
        .details { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; font-size: 14px; }
        .row { display: flex; justify-content: space-between; border-bottom: 1px dotted #ccc; padding-bottom: 5px; }
        .row span:first-child { color: #666; }
        .row span:last-child { font-weight: bold; }

        /* Footer / Button */
        .actions { margin-top: 10px; }
        button {
            width: 100%;
            padding: 15px;
            font-family: inherit;
            font-weight: bold;
            text-transform: uppercase;
            background: black;
            color: white;
            border: none;
            cursor: pointer;
            transition: opacity 0.2s;
        }
        button:hover { opacity: 0.8; }

        .footer-text { text-align: center; font-size: 10px; color: #999; margin-top: 15px; }
    </style>
</head>
<body>

    <div class="card">
        <header>
            <h1>System Report</h1>
            <div class="subtitle">Go Matrix • v1.0</div>
        </header>

        <div class="hero">
            <div class="hero-label">Performance Speedup</div>
            <div class="hero-val">{{printf "%.2f" .Speedup}}x</div>
            <div class="hero-note">Concurrent processing was {{printf "%.1f" .Speedup}} times faster</div>
        </div>

        <div class="viz-container">
            <div class="bar-group">
                <div class="bar-label">Sequential</div>
                <div class="track">
                    <div class="fill" style="width: 100%;"></div>
                </div>
                <div class="time-text">{{printf "%.2f" .SeqTime}}s</div>
            </div>

            <div class="bar-group">
                <div class="bar-label">Concurrent</div>
                <div class="track">
                    <div class="fill" style="width: {{printf "%.2f" .InverseSpeedup}}%;"></div>
                </div>
                <div class="time-text">{{printf "%.2f" .ConcTime}}s</div>
            </div>
        </div>

        <div class="details">
            <div class="row"><span>Matrix Size</span><span>{{.TotalRows}} x {{.TotalRows}}</span></div>
            <div class="row"><span>Total Elements</span><span>{{.ElementsFormatted}}</span></div>
            <div class="row"><span>Checksum</span><span>{{printf "%.2e" .Checksum}}</span></div>
            <div class="row"><span>Status</span><span>SUCCESS</span></div>
        </div>

				<div class="actions">
    				<button onclick="downloadCSV()">Download Report .CSV</button>
    				<div class="footer-text">
        			Created by M.Y & M.M • <a href="https://github.com/mohamadyoussefio/go-matrix-service" target="_blank" style="color: #666; text-decoration: underline;">View Source Code</a>
    				</div>
				</div>
    </div>

    <script>
        function downloadCSV() {
            const rows = [
                ["METRIC", "VALUE"],
                ["Speedup", "{{printf "%.2f" .Speedup}}x"],
                ["Sequential Time", "{{printf "%.4f" .SeqTime}}s"],
                ["Concurrent Time", "{{printf "%.4f" .ConcTime}}s"],
                ["Matrix Size", "{{.TotalRows}}"],
                ["Checksum", "{{printf "%.2e" .Checksum}}"]
            ];
            let csvContent = "data:text/csv;charset=utf-8," + rows.map(e => e.join(",")).join("\n");
            const encodedUri = encodeURI(csvContent);
            const link = document.createElement("a");
            link.setAttribute("href", encodedUri);
            link.setAttribute("download", "matrix_report.csv");
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        }
    </script>
</body>
</html>
`

func GenerateHTML(r protocol.Response) {
	type PageData struct {
		protocol.Response
		InverseSpeedup    float64
		ElementsFormatted string
		TimeStamp         string
	}

	data := PageData{
		Response:          r,
		InverseSpeedup:    (r.ConcTime / r.SeqTime) * 100.0,
		ElementsFormatted: fmt.Sprintf("%dM", (r.TotalRows*r.TotalRows)/1_000_000), // e.g. "4M"
		TimeStamp:         "2024-SESSION",
	}

	f, err := os.Create("report.html")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tmpl.Execute(f, data)
	fmt.Printf("\n\033[36m[REPORT] Generated clean report: 'report.html'\033[0m\n")
	openBrowser("report.html")
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
	}
}
