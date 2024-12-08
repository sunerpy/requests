window.BENCHMARK_DATA = {
  "lastUpdate": 1733635555690,
  "repoUrl": "https://github.com/sunerpy/requests",
  "entries": {
    "Requests-Benchmark": [
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "11db1ffcf1d4648b3c64e8a4b7f9e74e23cf9c62",
          "message": "ci(benchmark): 更新性能基准测试工作流\n\n- 使用 benchmark-action 生成和存储基准测试结果\n- 自动推送基准测试报告到 GitHub Pages",
          "timestamp": "2024-12-08T12:35:12+08:00",
          "tree_id": "58c433b8892c0c60893cd67a0b3fda08110a6947",
          "url": "https://github.com/sunerpy/requests/commit/11db1ffcf1d4648b3c64e8a4b7f9e74e23cf9c62"
        },
        "date": 1733632581001,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97402,
            "unit": "ns/op",
            "extra": "12364 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100970,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100899,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100560,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25397,
            "unit": "ns/op",
            "extra": "47257 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20508,
            "unit": "ns/op",
            "extra": "58208 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 24693,
            "unit": "ns/op",
            "extra": "58712 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20981,
            "unit": "ns/op",
            "extra": "57578 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "1f5a5f35c12639fc753627e55a5ed1e395ae395e",
          "message": "docs: 更新 README 文件并添加中文版本\n\n- 更新英文 README.md 文件内容\n- 新增中文 readme-cn.md 文件\n- 优化文档结构，增加目录和徽章\n- 更新示例代码和 API 文档",
          "timestamp": "2024-12-08T13:24:58+08:00",
          "tree_id": "6ca69b4379cc17bcd5eca42c26ece664d5426412",
          "url": "https://github.com/sunerpy/requests/commit/1f5a5f35c12639fc753627e55a5ed1e395ae395e"
        },
        "date": 1733635555406,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94744,
            "unit": "ns/op",
            "extra": "12327 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98135,
            "unit": "ns/op",
            "extra": "12172 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98580,
            "unit": "ns/op",
            "extra": "12192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 97966,
            "unit": "ns/op",
            "extra": "12200 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 27505,
            "unit": "ns/op",
            "extra": "48156 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20632,
            "unit": "ns/op",
            "extra": "55200 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20624,
            "unit": "ns/op",
            "extra": "58003 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20922,
            "unit": "ns/op",
            "extra": "58866 times\n4 procs"
          }
        ]
      }
    ]
  }
}