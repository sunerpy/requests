window.BENCHMARK_DATA = {
  "lastUpdate": 1733724732314,
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
          "id": "1389d0675c25752a642efe3278fbbdc0ea0b9fc8",
          "message": "feat(requests): 添加自定义 DNS 服务器支持和session重置\n\n- 在 Session 类型中增加 WithDNS 方法，支持设置自定义 DNS 服务器\n- 实现自定义 DNS 解析逻辑\n\n- 增加相关单元测试，验证自定义 DNS 功能的正确性\n\n- 增加了 Clear() 方法，确保会话状态重置",
          "timestamp": "2024-12-08T21:02:49+08:00",
          "tree_id": "848bb5b4250328645ac166e78ef396050097e0bc",
          "url": "https://github.com/sunerpy/requests/commit/1389d0675c25752a642efe3278fbbdc0ea0b9fc8"
        },
        "date": 1733663076655,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94821,
            "unit": "ns/op",
            "extra": "12531 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 97851,
            "unit": "ns/op",
            "extra": "12192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98365,
            "unit": "ns/op",
            "extra": "12158 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98077,
            "unit": "ns/op",
            "extra": "12205 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25751,
            "unit": "ns/op",
            "extra": "47043 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 23659,
            "unit": "ns/op",
            "extra": "45820 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20343,
            "unit": "ns/op",
            "extra": "57511 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20751,
            "unit": "ns/op",
            "extra": "59461 times\n4 procs"
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
          "id": "b4f9c20a400738d5719449b6bf8979608c50c3df",
          "message": "ci: 更新 GitHub Actions 工作流触发条件\n\n- 在 benchmark.yml 和 release.yml 中添加 dev 分支的推送触发\n- 在 benchmark.yml 中添加 main 分支的 PR 触发",
          "timestamp": "2024-12-09T11:45:54+08:00",
          "tree_id": "b1fd00253b2a1bcba1edde328c7e71cc62f2ab82",
          "url": "https://github.com/sunerpy/requests/commit/b4f9c20a400738d5719449b6bf8979608c50c3df"
        },
        "date": 1733716023147,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94034,
            "unit": "ns/op",
            "extra": "12760 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98214,
            "unit": "ns/op",
            "extra": "12235 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 99346,
            "unit": "ns/op",
            "extra": "12130 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98590,
            "unit": "ns/op",
            "extra": "12201 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 30015,
            "unit": "ns/op",
            "extra": "48313 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20370,
            "unit": "ns/op",
            "extra": "54230 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20655,
            "unit": "ns/op",
            "extra": "58539 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20634,
            "unit": "ns/op",
            "extra": "59032 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "b4f9c20a400738d5719449b6bf8979608c50c3df",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/b4f9c20a400738d5719449b6bf8979608c50c3df"
        },
        "date": 1733716157995,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95758,
            "unit": "ns/op",
            "extra": "12364 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100070,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100494,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 99559,
            "unit": "ns/op",
            "extra": "12010 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25265,
            "unit": "ns/op",
            "extra": "46749 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 23320,
            "unit": "ns/op",
            "extra": "58154 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20375,
            "unit": "ns/op",
            "extra": "57363 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20264,
            "unit": "ns/op",
            "extra": "59478 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "1a67006cc63ded39651f252807934bb674e47885",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/1a67006cc63ded39651f252807934bb674e47885"
        },
        "date": 1733721919738,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94590,
            "unit": "ns/op",
            "extra": "12471 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98893,
            "unit": "ns/op",
            "extra": "12566 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 99422,
            "unit": "ns/op",
            "extra": "12141 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100184,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24943,
            "unit": "ns/op",
            "extra": "48746 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 21030,
            "unit": "ns/op",
            "extra": "59098 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20916,
            "unit": "ns/op",
            "extra": "58503 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20651,
            "unit": "ns/op",
            "extra": "56932 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "d5660fcabb1a70b2adf3eb6e540a34d18925d849",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/d5660fcabb1a70b2adf3eb6e540a34d18925d849"
        },
        "date": 1733722063494,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97384,
            "unit": "ns/op",
            "extra": "12174 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100355,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 101188,
            "unit": "ns/op",
            "extra": "12022 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 101007,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24926,
            "unit": "ns/op",
            "extra": "46803 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20095,
            "unit": "ns/op",
            "extra": "58850 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20420,
            "unit": "ns/op",
            "extra": "58843 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 21216,
            "unit": "ns/op",
            "extra": "59574 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "1bbd1d863a1a3ca50db42746b55919829b494c99",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/1bbd1d863a1a3ca50db42746b55919829b494c99"
        },
        "date": 1733724731906,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95481,
            "unit": "ns/op",
            "extra": "12428 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98892,
            "unit": "ns/op",
            "extra": "12115 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98527,
            "unit": "ns/op",
            "extra": "12123 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98979,
            "unit": "ns/op",
            "extra": "12097 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 27281,
            "unit": "ns/op",
            "extra": "47341 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 22143,
            "unit": "ns/op",
            "extra": "45915 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20683,
            "unit": "ns/op",
            "extra": "58839 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20113,
            "unit": "ns/op",
            "extra": "59924 times\n4 procs"
          }
        ]
      }
    ]
  }
}