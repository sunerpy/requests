window.BENCHMARK_DATA = {
  "lastUpdate": 1733662384881,
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
          "id": "7ea5ab731a6c5432699c5f151dbcd6cb70f0c1ea",
          "message": "docs: 更新 README 文件并添加中文版本\n\n- 更新英文 README.md 文件内容\n- 新增中文 readme-cn.md 文件\n- 优化文档结构，增加目录和徽章\n- 更新示例代码和 API 文档\n- 新增 Benchmark 工作流，运行性能测试生成 HTML 报告并上传到 GitHub Pages",
          "timestamp": "2024-12-08T13:28:58+08:00",
          "tree_id": "6ca69b4379cc17bcd5eca42c26ece664d5426412",
          "url": "https://github.com/sunerpy/requests/commit/7ea5ab731a6c5432699c5f151dbcd6cb70f0c1ea"
        },
        "date": 1733635869599,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95525,
            "unit": "ns/op",
            "extra": "12519 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99849,
            "unit": "ns/op",
            "extra": "12178 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98913,
            "unit": "ns/op",
            "extra": "12192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98642,
            "unit": "ns/op",
            "extra": "12213 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24988,
            "unit": "ns/op",
            "extra": "48128 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20210,
            "unit": "ns/op",
            "extra": "59992 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20913,
            "unit": "ns/op",
            "extra": "60295 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20705,
            "unit": "ns/op",
            "extra": "58940 times\n4 procs"
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
          "id": "c8e51a9caf2a6eff1cb52238a31a7a5f31b1b9c9",
          "message": "feat(requests): 添加自定义 DNS 服务器支持\n\n- 在 Session 类型中增加 WithDNS 方法，支持设置自定义 DNS 服务器\n- 实现自定义 DNS 解析逻辑\n\n- 增加相关单元测试，验证自定义 DNS 功能的正确性",
          "timestamp": "2024-12-08T20:52:17+08:00",
          "tree_id": "55f8c734827a3da5f7ef356e94402bbbb2a1c2d6",
          "url": "https://github.com/sunerpy/requests/commit/c8e51a9caf2a6eff1cb52238a31a7a5f31b1b9c9"
        },
        "date": 1733662384100,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95805,
            "unit": "ns/op",
            "extra": "12405 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99465,
            "unit": "ns/op",
            "extra": "12058 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100949,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100439,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25286,
            "unit": "ns/op",
            "extra": "46009 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20489,
            "unit": "ns/op",
            "extra": "58611 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20552,
            "unit": "ns/op",
            "extra": "58785 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20754,
            "unit": "ns/op",
            "extra": "59725 times\n4 procs"
          }
        ]
      }
    ]
  }
}